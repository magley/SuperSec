# XSS

Објашњење: Апликација пропушта податке са једне странице без валидације на такав начин да је могуће уметнути нпр. JS који се може покренути.

Утицај: Крађа сесије, модификација сајтова, редирекција на злонамерне сајтове.

Рањивости: Употреба старе библиотеке, мањак валидације.

Контрамере: Користити најновију верзију библиотека за форме, за UI компоненте, за template engine-e и за обраду података на серверу. Вршити escaping карактера.

# Sensitive Data Exposure

Објашњење: Нападач приступа осетљивим подацима за које није овлашћен.

Утицај: Крађа финансијских и приватних података, модификација истих, представљање као друга особа итд.

Рањивости: Лоша контрола приступа. Незаштићене сировине.

Контрамере: Боља (строжа) контрола приступа: deny by default, форсирање аутентификације, ауторизације и идентификације. Омогућити да подаци буду сакривени док их корисник ручно не открије. Асиметрична криптографија приликом складиштења и слања.

# Improper Input Validation

Објашњење: Нападач искоришћава мањак валидација улаза да би извршио команде или приступио ресурсима за које није овлашћен.

Утицај: Крађа и измена података, arbitrary code execution (измена понашања апликација, уметање мелвера).

Рањивости: Лоша библиотека, неискорићене могућности валидације.

Контрамере: Користити савремене библиотеке које валидирају улаз и на клијентској и серверској страни. Escape-овање карактера. Форсирање одређеног regex формата.

# Broken Access Control

Објашњење: Нападач приступа функционалностима и подацима за које није овлашћен.

Утицај: Нападач добија увид у осетљиве податке корисника, платне податке, статистику апликације. Може да извршава команде у име других корисника.

Рањивости: Мањак контроле приступа.

Контрамере: Deny by default. RBAC, ABAC итд. (у зависности од типа софтвера).

# Unvalidated Redirects

Објашњење: Нападач може лажирати редирекцију ка малициозном сајту тако што злоупотреби наивну редирекцију.

Утицај: Корисник ће отворити phishing сајт и његови подаци ће бити украдени.

Рањивости: Редирекција се ради на основу URL-a и URL параметара без сложенијих валидација.

Контрамере: Валидација URL-a, валидација параметара са додатним подацима, обавештење о редирекцији.

# Unvalidated Redirects

Објашњење: Нападач може лажирати редирекцију ка малициозном сајту тако што злоупотреби наивну редирекцију.

Утицај: Корисник ће отворити phishing сајт и његови подаци ће бити украдени.

Рањивости: Редирекција се ради на основу URL-a и URL параметара без сложенијих валидација.

Контрамере: Валидација URL-a, валидација параметара са додатним подацима, обавештење о редирекцији.

# Vulnerable Components

Објашњење: Неке компоненте у апликацији су подложне нападима због рањивости. Преко њих је могуће злоупотребити друге компоненте у апликацији које су саме по себи заштићене.

Утицај: Неовлашћени приступ подацима. Arbitrary code execution.

Рањивости: Употреба старих библиотека. Веровање страним компонентама.

Контрамере: Употреба новијих библиотека. Zero trust policy.

# Broken Authentication

Објашњење: Апликација је лоше имплементирала аутентификацију те нападач се лако може представити као други корисник.

Утицај: Крађа сесија, фајлова. Приступ лозинкама.

Рањивости: Употреба старих библиотека. Лош дизајн у аутентификацији (нестандардни).

Контрамере: Придржавање стандарда. Употреба нових библиотека. MFA.

# Security through Obscurity

Објашњење: Апликација се ослања на алгоритме за систем заштите. Ако нападач открије интерне механизме заштите, лако ће их превазићи.

Утицај: Уништење читавог заштитног механизма.

Рањивости: Ослањање на интерне механизме заштите уместо проверених алгоритама.

Контрамере: Користити конкретне алгоритме које је тешко пробити. Јавно објавити који се алгоритми користе (али не и сами кључеви итд.).

# Insecure Deserialization

Објашњење: Лошом десеријализацијом нападач искоришћава рањивост у систему да изврши произвољан код.

Утицај: Извршавање произвољног кода, инјекција мелвера, уништавање сервера.

Рањивости: Употреба лоших библиотека и формата. Примена лоших алгоритама десеријализације.

Контрамере: Примена нових библиотека. Употреба текстуалних формата. Десеријализација на основу краја стринга уместо унапред датог броја.

# Broken Anti Automation

Објашњење: Нападач аутоматизује захтеве јер сервер није увео било какав вид контроле.

Утицај: DoS, DDOS

Рањивости: Сервер нема rate limiting.

Контрамере: Увести rate limiting.

# Injection

Објашњење: Нападач инјектује команду коју база података (или сличан сервис) извршава.

Утицај: Заобилажење контроле у бази података. Приступ неовлашћеним подацима. Брисање података.

Рањивости: Параметризација упита на основу конкатенације стрингова.

Контрамере: Prepared queries и слични механизми у другим сервисима. Употреба бољих библиотека и алата.

# Security Misconfiguration

Објашњење: Конфигурација сервера или неке његове компоненте није добро подешена.

Утицај: Неовлашћен приступ сировинама и функционалностима сервера.

Рањивости: Употреба подразумеване конфигурације, застареле конфигурације, лоших пракси при конфигурацији.

Контрамере: Правилна конфигурација компоненти на основу препорука произвођача тих компоненти и општих безбедносних пракси. Редовно ажурирање конфигурације у складу са изменама система.

# Cryptographic Issues

Објашњење: Систем који користи слабе криптографске алгоритме је подложан нападима.

Утицај: Добављање лозинки и других тајни. Компромитовање над подацима и функционалношћу.

Рањивости: Употреба застарелих и превазиђених криптографских алгоритама који се данас не препоручују. Чување кључева у тексуталним фајловима.

Контрамере: Миграција на новије, безбедније криптографске алгоритме. Чување кључева у key store-овима.

# XXE

Објашњење: Добављање спољашњих ресурса (који нису за дату ситуацију предвиђени) преко XML захтева.

Утицај: Неовлашћен приступ фајловима. Скенирање унутрашњих портова. Извршавање кодова на серверу. DoS напади.

Рањивости: Употреба библиотека за XML које дозвољавају external entity-е.

Контрамере: Ажурирање XML библиотека. Искључивање могућности за рад са external entity-има. Валидација улаза.