#include "esp_camera.h"
#include <WiFi.h>
#include <PubSubClient.h>
#include "DHT.h"

#define DHTPIN 14    // modify to the pin we connected
#define DHTTYPE DHT21   // AM2301 
DHT dht(DHTPIN, DHTTYPE);

// Select camera model
//#define CAMERA_MODEL_WROVER_KIT
//#define CAMERA_MODEL_ESP_EYE
//#define CAMERA_MODEL_M5STACK_PSRAM
//#define CAMERA_MODEL_M5STACK_WIDE
#define CAMERA_MODEL_AI_THINKER

#include "camera_pins.h"

const char* ssid = "ssid";
const char* password = "password";
const char* mqtt_server="m20.cloudmqtt.com"; //your mqtt server ip

WiFiClient espClient;
PubSubClient client(espClient);

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived ");
  Serial.println(topic);
  
  int64_t fr_start = esp_timer_get_time();
  camera_fb_t * fb = esp_camera_fb_get();
  if (!fb) {
      Serial.println("Camera capture failed");
  } else {
      // very slow
      // client.publish_P("/espCam/capture", fb->buf, fb->len, false); 

      // x25 faster
      //  client.beginPublish("/espCam/capture", fb->len, false);
      //  size_t meison = 0;
      //  size_t bufferLeft = fb->len;
      //  static const size_t bufferSize = 4096;
      //  static uint8_t buffer[bufferSize] = {0xFF};
      //  while (bufferLeft) {
      //    size_t copy = (bufferLeft < bufferSize) ? bufferLeft : bufferSize;
      //    memcpy ( &buffer, &fb->buf[meison], copy );
      //    client.write(&buffer[0], copy);
      //    // or just client.write(&fb->buf[meison], copy);
      //    bufferLeft -= copy;
      //    meison += copy;
      //  }
      //  client.endPublish();

      // x35 faster
      client.beginPublish("/espCam/capture", fb->len, false);
      client.write(&fb->buf[0], fb->len);
      client.endPublish();
      
      esp_camera_fb_return(fb);
      int64_t fr_end = esp_timer_get_time();
      Serial.printf("Sent JPG: %uB %ums\n", (uint32_t)(fb->len), (uint32_t)((fr_end - fr_start)/1000));
  }
}

void setup() {
  Serial.begin(115200);
  Serial.setDebugOutput(true);
  Serial.println();

  camera_config_t config;
  config.ledc_channel = LEDC_CHANNEL_0;
  config.ledc_timer = LEDC_TIMER_0;
  config.pin_d0 = Y2_GPIO_NUM;
  config.pin_d1 = Y3_GPIO_NUM;
  config.pin_d2 = Y4_GPIO_NUM;
  config.pin_d3 = Y5_GPIO_NUM;
  config.pin_d4 = Y6_GPIO_NUM;
  config.pin_d5 = Y7_GPIO_NUM;
  config.pin_d6 = Y8_GPIO_NUM;
  config.pin_d7 = Y9_GPIO_NUM;
  config.pin_xclk = XCLK_GPIO_NUM;
  config.pin_pclk = PCLK_GPIO_NUM;
  config.pin_vsync = VSYNC_GPIO_NUM;
  config.pin_href = HREF_GPIO_NUM;
  config.pin_sscb_sda = SIOD_GPIO_NUM;
  config.pin_sscb_scl = SIOC_GPIO_NUM;
  config.pin_pwdn = PWDN_GPIO_NUM;
  config.pin_reset = RESET_GPIO_NUM;
  config.xclk_freq_hz = 20000000;
  config.pixel_format = PIXFORMAT_JPEG;
  config.frame_size = FRAMESIZE_SVGA;
  config.jpeg_quality = 10;
  config.fb_count = 1;

#if defined(CAMERA_MODEL_ESP_EYE)
  pinMode(13, INPUT_PULLUP);
  pinMode(14, INPUT_PULLUP);
#endif

  // camera init
  esp_err_t err = esp_camera_init(&config);
  if (err != ESP_OK) {
    Serial.printf("Camera init failed with error 0x%x", err);
    return;
  }

  sensor_t * s = esp_camera_sensor_get();
  //initial sensors are flipped vertically and colors are a bit saturated
  if (s->id.PID == OV3660_PID) {
    s->set_vflip(s, 1);//flip it back
    s->set_brightness(s, 1);//up the blightness just a bit
    s->set_saturation(s, -2);//lower the saturation
  }
  
#if defined(CAMERA_MODEL_M5STACK_WIDE)
  s->set_vflip(s, 1);
  s->set_hmirror(s, 1);
#endif

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.println("WiFi connected");

  client.setServer(mqtt_server, 19001);
  client.setCallback(callback);
  client.connect("ESPCAM","user","password");
  client.publish("/espCam/info","Started");
  client.subscribe("/espCam/get");
  
  dht.begin();
}

void recon() {
    while (!client.connected()) {
      if(client.connect("ESPCAM")){
        client.publish("/espCam/info","Reconnected");
      }else{
        delay(1000);
      }
    }
}

void dht_loop() {
  static long lastMsg = 0;
  long now = millis();
  if (now - lastMsg > 60000) {
    lastMsg = now;
    
    float t = dht.readTemperature();
    if (isnan(t)) {
      Serial.println("Could not get temperature"); 
    } else {
      Serial.print("Temperature: "); 
      Serial.print(t);
      Serial.println(" *C");
      client.publish("/am2301/temperature",String(t).c_str());
    }

    float h = dht.readHumidity();
    if (isnan(h)) {
      Serial.println("Could not get humidity"); 
    } else {
      Serial.print("Humidity: "); 
      Serial.print(h);
      Serial.println(" %");
      client.publish("/am2301/humidity",String(h).c_str());
    }
  }
}

void loop() {
  if (!client.connected()) {
    recon();
  }
  client.loop();
  dht_loop();
}
