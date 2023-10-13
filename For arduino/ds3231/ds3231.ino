// Date and time functions using a DS3231 RTC connected via I2C and Wire lib
#include "RTClib.h"

RTC_DS3231 rtc;

char daysOfTheWeek[7][12] = {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"};
const long interval = 1000; 
unsigned long previousMillis = 0;  // will store last time LED was updated

#define PARSE_AMOUNT 6         // число значений в массиве, который хотим получить
int intData[PARSE_AMOUNT];     // массив численных значений после парсинга
boolean recievedFlag;
boolean getStarted;
byte index;
String string_convert = "";

/*
   Данный алгоритм позволяет получить через Serial пачку значений, и раскидать
   их в целочисленный массив. Использовать можно банально для управления
   ЧЕМ УГОДНО через bluetooth, так как bluetooth модули есть UART интерфейс связи.
   Либо управлять через Serial с какой-то программы с ПК
   Как использовать:
   1) В PARSE_AMOUNT указывается, какое количество значений мы хотим принять.
   От этого значения напрямую зависит размер массива принятых данных, всё просто
   2) Пакет данных на приём должен иметь вид:
   Начало - символ $
   Разделитель - пробел
   Завершающий символ - ;
   Пример пакета: $110 25 600 920;  будет раскидан в массив intData согласно порядку слева направо
   Что делает данный скетч:
   Принимает пакет данных указанного выше вида, раскидывает его в массив intData, затем выводит обратно в порт.
   Отличие от предыдущего примера: написан мной, не используя никаких хитрых функций. Предельно просто и понятно работает
*/

void parsing() {
  if (Serial.available() > 0) {
    char incomingByte = Serial.read();        // обязательно ЧИТАЕМ входящий символ
    
    if (getStarted) {                         // если приняли начальный символ (парсинг разрешён)
      if (incomingByte != ' ' && incomingByte != ';') {   // если это не пробел И не конец
        string_convert += incomingByte;       // складываем в строку
      } else {                                // если это пробел или ; конец пакета
        intData[index] = string_convert.toInt();  // преобразуем строку в int и кладём в массив
        string_convert = "";                  // очищаем строку
        index++;                              // переходим к парсингу следующего элемента массива
      }
    }
    if (incomingByte == '$') {                // если это $
      getStarted = true;                      // поднимаем флаг, что можно парсить
      index = 0;                              // сбрасываем индекс
      string_convert = "";                    // очищаем строку
    }
    if (incomingByte == ';') {                // если таки приняли ; - конец парсинга
      getStarted = false;                     // сброс
      recievedFlag = true;                    // флаг на принятие
      Serial.print('@');
    }
  }
}

void setup () {
  Serial.begin(115200);

#ifndef ESP8266
  while (!Serial); // wait for serial port to connect. Needed for native USB
#endif

  if (! rtc.begin()) {
    Serial.println("Couldn't find RTC");
    Serial.flush();
    while (1) delay(10);
  }
  // rtc.adjust(DateTime(2014, 1, 21, 3, 0, 0));
}

void loop () {

    parsing();       // функция парсинга

    DateTime now = rtc.now();
    unsigned long currentMillis = millis();
    
    if (currentMillis - previousMillis >= interval) {
		previousMillis = currentMillis;
		Serial.print(now.year(), DEC);
		Serial.print('-');
		if(now.month()<10) {Serial.print("0");}
		Serial.print(now.month(), DEC);
		Serial.print('-');
		if(now.day()<10) {Serial.print("0");}
		Serial.print(now.day(), DEC);
		Serial.print(' ');
		if(now.hour()<10) {Serial.print("0");}
		Serial.print(now.hour(), DEC);
		Serial.print(':');
		if(now.minute()<10) {Serial.print("0");}   
		Serial.print(now.minute(), DEC);
		Serial.print(':');
		if(now.second()<10) {Serial.print("0");}  
		Serial.print(now.second(), DEC);
		Serial.println();
    }


  if (recievedFlag) {                           // если получены данные
    recievedFlag = false;
    rtc.adjust(DateTime(intData[0], intData[1], intData[2], intData[3], intData[4], intData[5]));
  }
}
