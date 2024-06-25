У овом документу налази се анализа OWASP ASVS 4.0.3 документа примењен на Lesotho пројекат.

## V1 Architecture, Design and Threat Modeling

### V1.1 Secure Software Development Lifecycle

Није разматрано. Овај пројекат је сувише малих размера и разбијање на целине,
корисничке приче итд.; а ради дефинисања безбедносних апсеката није исплативо.

### V1.2 Authentication Architecture

| Број | Опис | CWE | Испуњеност |
| ---- | ---- | --- | ---------- |
| 1.2.1| Verify the use of unique or special low-privilege operating system accounts for all application components, services, and servers.| 250 | **Да**, ниједан сервис не захтева администраторске привилегије. |
| 1.2.2| Verify that communications between application components, including APIs, middleware and data layers, are authenticated. Components should have the least necessary privileges needed. | 306 | **Парцијално**, API кључ је имплементиран али интерне компоненте су подешене са подразумеваним привилегијама. |
| 1.2.3| Verify that the application uses a single vetted authentication mechanism that is known to be secure, can be extended to include strong authentication, and has sufficient logging and monitoring to detect account abuse or breaches. | 306 | **Парцијално**, logging јесте, а monitoring није имплементиран. |
| 1.2.4| Verify that all authentication pathways and identity management APIs implement consistent authentication security control strength, such that there are no weaker alternatives per the risk of the application.| 306 | **Да**, имплементиран је само један аутентификациони (_demo2_) и само један ауторизациони (_централизовани Lesotho auth_) систем у склопу пројекта.
---

### V1.3 Session Management Architecture

> This is a placeholder for future architectural requirements.

### V1.4 Access Control Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.4.1 | Verify that trusted enforcement points, such as access control gateways, servers, and serverless functions, enforce access controls. Never enforce access controls on the client. | 602 | **Да** |
| 1.4.2 | [DELETED, NOT ACTIONABLE] |
| 1.4.3 | [DELETED, DUPLICATE OF 4.1.3] |
| 1.4.4 | Verify the application uses a single and well-vetted access control mechanism for accessing protected data and resources. All requests must pass through this single mechanism to avoid copy and paste or insecure alternative paths. | 284 | **Да**, коришћен је Lesotho. Demo2 шаље ауторизационе захтеве серверу при покушају читања, писања и дељења документа. |
| 1.4.5 | Verify that attribute or feature-based access control is used whereby the code checks the user's authorization for a feature/data item rather than just their role. Permissions should still be allocated using roles. | 275 | **Не**, коришћен је наивни RBAC где се пермисије успостављају на нивоу рола |

### V1.5 Input and Output Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.5.1 | Verify that input and output requirements clearly define how to handle and process data based on type, content, and applicable laws, regulations, and other policy compliance. | 1029 | **?**
| 1.5.2 | Verify that serialization is not used when communicating with untrusted clients. If this is not possible, ensure that adequate integrity controls (and possibly encryption if sensitive data is sent) are enforced to prevent deserialization attacks including object injection. | 502 | **Да**, namespace и ACL се валидирају на пријему. |
| 1.5.3 | Verify that input validation is enforced on a trusted service layer. | 602 | **Да**, све валидације се извршавају на серверу. Потврда лозинке се ради на клијентском слоју, међутим то се не сматра безбедносним пропустом. |
| 1.5.4 | Verify that output encoding occurs close to or by the interpreter for which it is intended. | 116 | **?** |

### V1.6 Cryptographic Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.6.1 | Verify that there is an explicit policy for management of cryptographic keys and that a cryptographic key lifecycle follows a key management standard such as NIST SP 800-57. | 320 | **N/A** |
| 1.6.2 | Verify that consumers of cryptographic services protect key material and other secrets by using key vaults or API based alternatives. | 320 | **N/A** |
| 1.6.3 | Verify that all keys and passwords are replaceable and are part of a well-defined process to re-encrypt sensitive data. | 320 | **Не**, шифре нису заменљиве. |
| 1.6.4 | Verify that the architecture treats client-side secrets--such as symmetric keys, passwords, or API tokens--as insecure and never uses them to protect or access sensitive data. | 320 | **Да**, шифре се користе искључиво за аутентификацију. |

### V1.7 Errors, Logging and Auditing Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.7.1 | Verify that a common logging format and approach is used across the system. | 1009 | **Парцијално**, Lesotho има један формат, а Flask сервери други. |
| 1.7.2 | Verify that logs are securely transmitted to a preferably remote system for analysis, detection, alerting, and escalation. | | **Не**, све је на једном рачунару. |

### V1.8 Data Protection and Privacy Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.8.1 | Verify that all sensitive data is identified and classified into protection levels. | | **Не**, сви су у истом нивоу заштите (мрежном). |
| 1.8.2 | 1.8.2 Verify that all protection levels have an associated set of protection requirements, such as encryption requirements, integrity requirements, retention, privacy and other confidentiality requirements, and that these are applied in the architecture. | | **Не** |

### V1.9 Communications Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.9.1 | Verify the application encrypts communications between components, particularly when these components are in different containers, systems, sites, or cloud providers. | 319 | **Не**, HTTPs није имплементиран. |
| 1.9.2 | Verify that application components verify the authenticity of each side in a communication link to prevent person-in-the-middle attacks. For example, application components should validate TLS certificates and chains. | 295 | **Не**, HTTPs није имплементиран. |

### V1.10 Malicious Software Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.10.1 | Verify that a source code control system is in use, with procedures to ensure that check-ins are accompanied by issues or change tickets. The source code control system should have access control and identifiable users to allow traceability of any changes. | 284 | **Да**, коришћен је git, репозиторијум је на github серверу, и _issues_ су **делимично коришћени**. |

### V1.11 Business Logic Architecture

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 1.11.1 | Verify the definition and documentation of all application components in terms of the business or security functions they provide. | 1059 | **Да**, пословна логика функционалности, њихова документација и јединични тестови су усаглашћени.
| 1.11.2 | Verify that all high-value business logic flows, including authentication, session management and access control, do not share unsynchronized state. | 362 | **N/A**, над једном инстанцом је захтев испуњен. Скалирање сервиса није тестирано.
| 1.11.3 | Verify that all high-value business logic flows, including authentication, session management and access control are thread safe and resistant to time-of-check and time-of-use race conditions. | 367 | **? (Не)**, LevelDB и Consul-ов KV store су _thread-safe_, док Lesotho-ова кеш меморија није тестирана али вероватно није _thread-safe_. |

### V1.12 Secure File Upload Architecture

| Број   | Опис | CWE | Испуњеност |
| ------ | ---- | --- | ---------- |
| 1.12.1 | [DELETED, DUPLICATE OF 12.4.1] | |
| 1.12.2 | Verify that user-uploaded files - if required to be displayed or downloaded from the application - are served by either octet stream downloads, or from an unrelated domain, such as a cloud file storage bucket. Implement a suitable Content Security Policy (CSP) to reduce the risk from XSS vectors or other attacks from the uploaded file. | 646 | **Да**, фајлови се не приказују нити преузимају. |

### V1.13 API Architecture

> This is a placeholder for future architectural requirements.

### V1.14 Configuration Architecture

| Број   | Опис | CWE | Испуњеност |
| ------ | ---- | --- | ---------- |
| 1.14.1 | Verify the segregation of components of differing trust levels through well-defined security controls, firewall rules, API gateways, reverse proxies, cloud-based security groups, or similar mechanisms. | 923 | **N/A**, читав пројекат се извршава локално.
| 1.14.2 | Verify that binary signatures, trusted connections, and verified endpoints are used to deploy binaries to remote devices. | 494 | **N/A** |
| 1.14.3 | Verify that the build pipeline warns of out-of-date or insecure components and takes appropriate actions. | 1104 | **N/A** |
| 1.14.4 | Verify that the build pipeline contains a build step to automatically build and verify the secure deployment of the application, particularly if the application infrastructure is software defined, such as cloud environment build scripts. | | **Не**, рађено је ручно. |
| 1.14.5 | Verify that application deployments adequately sandbox, containerize and/or isolate at the network level to delay and deter attackers from attacking other applications, especially when they are performing sensitive or dangerous actions such as deserialization. | 265 | **Не** |
| 1.14.6 | Verify the application does not use unsupported, insecure, or deprecated client-side technologies such as NSAPI plugins, Flash, Shockwave, ActiveX, Silverlight, NACL, or client-side Java applets. | 477 | **Да** |

## V2 Authentication

### V2.1 Password Security

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.1.1 | Verify that user set passwords are at least 12 characters in length (after multiple spaces are combined) | 521 | 5.1.1.2 | **Да** |
| 2.1.2 | Verify that passwords of at least 64 characters are permitted, and that passwords of more than 128 characters are denied. | 521 | 5.1.1.2 | **Да** |
| 2.1.3 | Verify that password truncation is not performed. However, consecutive multiple spaces may be replaced by a single space. | 521 | 5.1.1.2 | **Да** |
| 2.1.4 | Verify that any printable Unicode character, including language neutral characters such as spaces and Emojis are permitted in passwords. | 521 | 5.1.1.2 | **Не**, коришћен је `input()`. |
| 2.1.5 | Verify users can change their password. | 620 | 5.1.1.2 | **Не** |
| 2.1.6 | Verify that password change functionality requires the user's current and new password. | 620 | 5.1.1.2 | **N/A** |
| 2.1.7 | Verify that passwords submitted during account registration, login, and password change are checked against a set of breached passwords either locally (such as the top 1,000 or 10,000 most common passwords which match the system's password policy) or using an external API. If using an API a zero knowledge proof or other mechanism should be used to ensure that the plain text password is not sent or used in verifying the breach status of the  password. If the password is breached, the application must require the user to set a new non-breached password. | 521 | 5.1.1.2 | **Не** |
| 2.1.8 | Verify that a password strength meter is provided to help users set  a stronger password. | 521 | 5.1.1.2 | **Не**, јер је коришћен CLI |
| 2.1.9 | Verify that there are no password composition rules limiting the  type of characters permitted. There should be no requirement for  upper or lower case or numbers or special characters. | 521 | 5.1.1.2 | **Да** |
| 2.1.10 | Verify that there are no periodic credential rotation or password  history requirements. | 263 | 5.1.1.2 | **Да** |
| 2.1.11 | Verify that "paste" functionality, browser password helpers, and  external password managers are permitted. | 521 | 5.1.1.2 | **Да** |
| 2.1.12 | Verify that the user can choose to either temporarily view the  entire masked password, or temporarily view the last typed character of the password on platforms that do not have this as built-in functionality. | 521 | 5.1.1.2 | **Не** |

### V2.2 General Authenticator Security

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.2.1 | Verify that anti-automation controls are effective at mitigating  breached credential testing, brute force, and account lockoutattacks. Such controls include blocking the most common  breached passwords, soft lockouts, rate limiting, CAPTCHA, ever  increasing delays between attempts, IP address restrictions, or  risk-based restrictions such as location, first login on a device,  recent attempts to unlock the account, or similar. Verify that no  more than 100 failed attempts per hour is possible on a single  account. | 307 | 5.2.2 / 5.1.1.2 / 5.1.4.2 / 5.1.5.2 | **Не** |
| 2.2.2 | Verify that the use of weak authenticators (such as SMS and email) is limited to secondary verification and transaction approval and not as a replacement for more secure authentication methods. Verify that stronger methods are offered before weak methods, users are aware of the risks, or that proper measures are in place to limit the risks of account compromise. | 304 | 5.2.10.0 | **Не** |
| 2.2.3 | Verify that secure notifications are sent to users after updates to authentication details, such as credential resets, email or address changes, logging in from unknown or risky locations. The use of push notifications - rather than SMS or email - is preferred, but in the absence of push notifications, SMS or email is acceptable as long as no sensitive information is disclosed in the notification. | 620 | | **Не** |
| 2.2.4 | Verify impersonation resistance against phishing, such as the use of multi-factor authentication, cryptographic devices with intent (such as connected keys with a push to authenticate), or at higher AAL levels, client-side certificates. | 308 | 5.2.5 | **Не**
| 2.2.5 | Verify that where a Credential Service Provider (CSP) and the  application verifying authentication are separated, mutually  authenticated TLS is in place between the two endpoints. | 319 | 5.2.6 | **Не** |
| 2.2.6 | Verify replay resistance through the mandated use of One-time Passwords (OTP) devices, cryptographic authenticators, or lookup codes. | 308 | 5.2.8 | **Не** |
| 2.2.7 | Verify intent to authenticate by requiring the entry of an OTP token or user-initiated action such as a button press on a FIDO hardware key. | 308 | 5.2.9 | **Не** |

### V2.3 Authenticator Lifecycle

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.3.1 | Verify system generated initial passwords or activation codes  SHOULD be securely randomly generated, SHOULD be at least 6 characters long, and MAY contain letters and numbers, and expire after a short period of time. These initial secrets must not be permitted to become the long term password. | 330 | 5.1.1.2 | **N/A** |
| 2.3.2 | Verify that enrollment and use of user-provided authentication  devices are supported, such as a U2F or FIDO tokens. | 308 | 6.1.3 | **Не**|
| 2.3.3 | Verify that renewal instructions are sent with sufficient time to  renew time bound authenticators. | 287 | 6.1.4 | **Не** |

### V2.4 Credential Storage

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.4.1 | Verify that passwords are stored in a form that is resistant to offline attacks. Passwords SHALL be salted and hashed using an approved one-way key derivation or password hashing function. Key derivation and password hashing functions take a password, a salt, and a cost factor as inputs when generating a password hash. | 916 | 5.1.1.2 | **Да**, коришћен је _bcrypt_. |
| 2.4.2 | Verify that the salt is at least 32 bits in length and be chosen arbitrarily to minimize salt value collisions among stored hashes. For each credential, a unique salt value and the resulting hash SHALL be stored. | 916 | 5.1.1.2 | **Да** |
| 2.4.3 | Verify that if PBKDF2 is used, the iteration count SHOULD be as large as verification server performance will allow, typically at least 100,000 iterations. | 916 | 5.1.1.2 | **N/A** |
| 2.4.4 | Verify that if bcrypt is used, the work factor SHOULD be as large as  verification server performance will allow, with a minimum of 10. | 916 | 5.1.1.2 | **Да**, коришћена је подразумевана вредност од 12 ([извор](https://pypi.org/project/bcrypt/#:~:text=Adjustable%20Work%20Factor)) |
| 2.4.5 | Verify that an additional iteration of a key derivation function is performed, using a salt value that is secret and known only to the verifier. Generate the salt value using an approved random bit generator [SP 800-90Ar1] and provide at least the minimum security strength specified in the latest revision of SP 800-131A. The secret salt value SHALL be stored separately from the hashed passwords (e.g., in a specialized device like a hardware security module). | 916 | 5.1.1.2 | **Не** |

### V2.5 Credential Recovery

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 2.5.1 | Verify that a system generated initial activation or recovery secret is not sent in clear text to the user. | 640 | 5.1.1.2 | **N/A** |
| 2.5.2 | Verify password hints or knowledge-based authentication (so-called "secret questions") are not present | 640 | 5.1.1.2 | **Да** |
| 2.5.3 | Verify password credential recovery does not reveal the current password in any way. | 640 | 5.1.1.2 | **N/A** |
| 2.5.4 | Verify shared or default accounts are not present (e.g. "root", "admin", or "sa"). | 16 | 5.1.1.2 | **Да** |
| 2.5.5 | Verify that if an authentication factor is changed or replaced, that the user is notified of this event. | 304 | 6.1.2.3 | **N/A** |
| 2.5.6 | Verify forgotten password, and other recovery paths use a secure recovery mechanism, such as time-based OTP (TOTP) or other soft token, mobile push, or another offline recovery mechanism. | 640 | 5.1.1.2 | **N/A** |
| 2.5.7 | Verify that if OTP or multi-factor authentication factors are lost, that evidence of identity proofing is performed at the same level as during enrollment. | 308 | 6.1.2.3 | **N/A** |

### V2.6 Look-up Secret Verifier

Није апликабилно за овај пројекат.

### V2.7 Out of Band Verifier

Није апликабилно за овај пројекат.

### V2.8 One Time Verifier

Није апликабилно за овај пројекат.

### V2.9 Cryptographic Verifier

Није апликабилно за овај пројекат.

### V2.10 Service Authentication

Није апликабилно за овај пројекат.

## V3 Session Management

### V3.1 Fundamental Session Management Security

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 3.1.1 | Verify the application never reveals session tokens in URL parameters. | 598 | | **Да** |

### V3.2 Session Binding

Није апликабилно за овај пројекат.

### V3.3 Session Termination

Није апликабилно за овај пројекат. Сесије су локалне и чувају се у радној меморији клијента (demo2).

### V3.4 Cookie-based Session Management

Није апликабилно за овај пројекат. Не користе се кукији.

### V3.5 Cookie-based Session Management

Није апликабилно за овај пројекат. Не користе се JWT/OAuth/SAML кључеви.

### V3.6 Federated Re-authentication

Није апликабилно за овај пројекат.

### V3.7 Defenses Against Session Management Exploits

| Број  | Опис | CWE | NIST | Испуњеност |
| ----- | ---- | --- | ---- | ---------- |
| 3.7.1 | Verify the application ensures a full, valid login session or requires re-authentication or secondary verification before allowing any sensitive transactions or account modifications. | 306 | | **Не** |

## V4 Access Control

### V4.1 General Access Control Design

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 4.1.1 | Verify that the application enforces access control rules on a trusted service layer, especially if client-side access control is present and could be bypassed. | 602 | **Да** |
| 4.1.2 | Verify that all user and data attributes and policy information used by access controls cannot be manipulated by end users unless specifically authorized. | 639 | **Да** |
| 4.1.3 | Verify that the principle of least privilege exists - users should only be able to access functions, data files, URLs, controllers, services, and other resources, for which they possess specific authorization. This implies protection against spoofing and elevation of privilege.| 285 | **?** |
| 4.1.4 | [DELETED, DUPLICATE OF 4.1.3] | | |
| 4.1.5 | Verify that access controls fail securely including when an exception occurs. | 285 | **Не** | 

### V4.2 Operation Level Access Control

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 4.2.1 | Verify that sensitive data and APIs are protected against Insecure Direct Object Reference (IDOR) attacks targeting creation, reading, updating and deletion of records, such as creating or updating someone else's record, viewing everyone's records, or deleting all records. | 639 | **Да**, сваком захтеву претходи access control. |
| 4.2.2 | Verify that the application or framework enforces a strong anti-CSRF mechanism to protect authenticated functionality, and effective anti-automation or anti-CSRF protects unauthenticated functionalit | 352 | **Да**, CSRF није могућ. |

### V4.3 Other Access Control Considerations

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 4.3.1 | Verify administrative interfaces use appropriate multi-factor authentication to prevent unauthorized use. | 419 | **N/A** |
| 4.3.2 | Verify that directory browsing is disabled unless deliberately desired. Additionally, applications should not allow discovery or disclosure of file or directory metadata, such as Thumbs.db, .DS_Store, .git or .svn folders. | 548 | **N/A** |
| 4.3.3 | Verify the application has additional authorization (such as step up or adaptive authentication) for lower value systems, and / or segregation of duties for high value applications to enforce anti-fraud controls as per the risk of application and past fraud. | 732 | **N/A** |

## V5 Validation, Sanitization and Encoding

### V5.1 Input Validation

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.1.1 | Verify that the application has defenses against HTTP parameter pollution  attacks, particularly if the application framework makes no distinction about  the source of request parameters (GET, POST, cookies, headers, or  environment variables). | 235 | **Да** |
| 5.1.2 | Verify that frameworks protect against mass parameter assignment attacks,  or that the application has countermeasures to protect against unsafe  parameter assignment, such as marking fields private or similar. | 915 | **Да** |
| 5.1.3 | Verify that all input (HTML form fields, REST requests, URL parameters, HTTP  headers, cookies, batch files, RSS feeds, etc) is validated using positive  validation (allow lists). | 20 | **Не** |
| 5.1.4 | Verify that structured data is strongly typed and validated against a defined  schema including allowed characters, length and pattern (e.g. credit card  numbers, e-mail addresses, telephone numbers, or validating that two related  fields are reasonable, such as checking that suburb and zip/postcode match). | 20 | **Не** (парцијална валидација постоји) |
| 5.1.5 | Verify that URL redirects and forwards only allow destinations which appear  on an allow list, or show a warning when redirecting to potentially untrusted  content. | 601 | **N/A** |

### V5.2 Sanitization and Sandboxing

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.2.1 | Verify that all untrusted HTML input from WYSIWYG editors or similar is  properly sanitized with an HTML sanitizer library or framework feature. | 116 | **N/A** |
| 5.2.2 | Verify that unstructured data is sanitized to enforce safety measures such as  allowed characters and length. | 138 | **Не** |
| 5.2.3 | Verify that the application sanitizes user input before passing to mail systems  to protect against SMTP or IMAP injection. | 147 | **N/A** |
| 5.2.4 | Verify that the application avoids the use of eval() or other dynamic code  execution features. Where there is no alternative, any user input being  included must be sanitized or sandboxed before being executed. | 95 | **Да** |
| 5.2.5 | Verify that the application protects against template injection attacks by  ensuring that any user input being included is sanitized or sandboxed. | 94 | **N/A** |
| 5.2.6 | Verify that the application protects against SSRF attacks, by validating or  sanitizing untrusted data or HTTP file metadata, such as filenames and URL  input fields, and uses allow lists of protocols, domains, paths and ports. | 918 | **N/A** |
| 5.2.7 | Verify that the application sanitizes, disables, or sandboxes user-supplied  Scalable Vector Graphics (SVG) scriptable content, especially as they relate to  XSS resulting from inline scripts, and foreignObject. | 159 | **N/A** |
| 5.2.8 | Verify that the application sanitizes, disables, or sandboxes user-supplied  scriptable or expression template language content, such as Markdown, CSS  or XSL stylesheets, BBCode, or similar. | 94 | **N/A** |

### V5.3 Output Encoding and Injection Prevention

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.3.1 | Verify that output encoding is relevant for the interpreter and context  required. For example, use encoders specifically for HTML values, HTML  attributes, JavaScript, URL parameters, HTTP headers, SMTP, and others as  the context requires, especially from untrusted inputs (e.g. names with  Unicode or apostrophes, such as ねこ or O'Hara). | 116 | **Не** |
| 5.3.2 | Verify that output encoding preserves the user's chosen character set and  locale, such that any Unicode character point is valid and safely handled. | 176 | **Не** |
| 5.3.3 | Verify that context-aware, preferably automated - or at worst, manual - output escaping protects against reflected, stored, and DOM based XSS. | 79 | **Не** |
| 5.3.4 | Verify that data selection or database queries (e.g. SQL, HQL, ORM, NoSQL)  use parameterized queries, ORMs, entity frameworks, or are otherwise  protected from database injection attacks. | 89 | **N/A** |
| 5.3.5 | Verify that where parameterized or safer mechanisms are not present,  context-specific output encoding is used to protect against injection attacks,  such as the use of SQL escaping to protect against SQL injection. | 89 | **N/A** |
| 5.3.6 | Verify that the application protects against JSON injection attacks, JSON eval  attacks, and JavaScript expression evaluation. | 830 | **?** |
| 5.3.7 | Verify that the application protects against LDAP injection vulnerabilities, or  that specific security controls to prevent LDAP injection have been  implemented. | 90 | **N/A** |
| 5.3.8 | Verify that the application protects against OS command injection and that  operating system calls use parameterized OS queries or use contextual  command line output encoding. | 78 | **N/A** |
| 5.3.9 | Verify that the application protects against Local File Inclusion (LFI) or  Remote File Inclusion (RFI) attacks. | 829 | **N/A** |
| 5.3.10 | Verify that the application protects against XPath injection or XML injection  attacks.  | 643 | **N/A** |

### V5.4 Memory, String, and Unmanaged Code

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.4.1 | Verify that the application uses memory-safe string, safer memory copy and  pointer arithmetic to detect or prevent stack, buffer, or heap overflows. | 120 | **?** |
| 5.4.2 | Verify that format strings do not take potentially hostile input, and are  constant. | 134 | **Не** (парцијално задовољено) |
| 5.4.3 | Verify that sign, range, and input validation techniques are used to prevent  integer overflows. | 190 | **N/A** |

### V5.5 Deserialization Prevention

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 5.5.1 | Verify that serialized objects use integrity checks or are encrypted to prevent  hostile object creation or data tampering. | 502 | **Да** |
| 5.5.2 | Verify that the application correctly restricts XML parsers to only use the most  restrictive configuration possible and to ensure that unsafe features such as  resolving external entities are disabled to prevent XML eXternal Entity (XXE)  attacks. | 611 | **N/A** |
| 5.5.3 | Verify that deserialization of untrusted data is avoided or is protected in both  custom code and third-party libraries (such as JSON, XML and YAML parsers). | 502 | **Не** |
| 5.5.4 | Verify that when parsing JSON in browsers or JavaScript-based backends,  JSON.parse is used to parse the JSON document. Do not use eval() to parse  JSON. | 95 | **Да** |

## V6 Stored Cryptography

Ово поглавље није разматрано за овај пројекат.

## V7 Error Handling and Logging

### V7.1 Log Content

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.1.1 | Verify that the application does not log credentials or payment details.  Session tokens should only be stored in logs in an irreversible, hashed form. | 532 | **Да** |
| 7.1.2 | Verify that the application does not log other sensitive data as defined under  local privacy laws or relevant security policy. | 532 | **Да** |
| 7.1.3 | Verify that the application logs security relevant events including successful  and failed authentication events, access control failures, deserialization  failures and input validation failures. | 778 | **Да** |
| 7.1.4 | Verify that each log event includes necessary information that would allow for  a detailed investigation of the timeline when an event happens. | 778 | **Да** |

### V7.2 Log Processing

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.2.1 | Verify that all authentication decisions are logged, without storing sensitive  session tokens or passwords. This should include requests with relevant  metadata needed for security investigations. | 778 | **Да** |
| 7.2.2 | Verify that all access control decisions can be logged and all failed decisions  are logged. This should include requests with relevant metadata needed for  security investigations. | 285 | **Да** |

### V7.3 Log Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.3.1 | Verify that all logging components appropriately encode data to prevent log  injection. | 117 | **N/A** |
| 7.3.2 | [DELETED, DUPLICATE OF 7.3.1] |  | |
| 7.3.3 | Verify that security logs are protected from unauthorized access and  modification. | 200 | **N/A** |
| 7.3.4 | Verify that time sources are synchronized to the correct time and time zone.  Strongly consider logging only in UTC if systems are global to assist with post-incident forensic analysis. |  | **Не** |

### V7.4 Error Handling

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 7.4.1 | Verify that a generic message is shown when an unexpected or security  sensitive error occurs, potentially with a unique ID which support personnel  can use to investigate.  | 210 | **Не**, порука је дескриптивна |
| 7.4.2 | Verify that exception handling (or a functional equivalent) is used across the  codebase to account for expected and unexpected error conditions. | 533 | **Да** |
| 7.4.3 | Verify that a "last resort" error handler is defined which will catch all  unhandled exceptions.  | 431 | **?**, није познато за сам Lesotho сервис где је коришћен Go. |

## V8 Data Protection

### V8.1 General Data Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 8.1.1 | Verify the application protects sensitive data from being cached in server  components such as load balancers and application caches. | 524 | **Да**, кеширање је ручно и једино се кешира namespace који нема осетљиве податке. |
| 8.1.2 | Verify that all cached or temporary copies of sensitive data stored on the  server are protected from unauthorized access or purged/invalidated after  the authorized user accesses the sensitive data. | 524 | **N/A**, |
| 8.1.3 | Verify the application minimizes the number of parameters in a request, such  as hidden fields, Ajax variables, cookies and header values. | 233 | **Да** |
| 8.1.4 | Verify the application can detect and alert on abnormal numbers of requests,  such as by IP, user, total per hour or day, or whatever makes sense for the  application. | 770 | **Не** |
| 8.1.5 | Verify that regular backups of important data are performed and that test  restoration of data is performed. | 19 | **Не** |
| 8.1.6 | Verify that backups are stored securely to prevent data from being stolen or  corrupted. | 19 | **N/A** |

### V8.2 Client-side Data Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 8.2.1 | Verify the application sets sufficient anti-caching headers so that sensitive  data is not cached in modern browsers. | 525 | **N/A** |
| 8.2.2 | Verify that data stored in browser storage (such as localStorage,  sessionStorage, IndexedDB, or cookies) does not contain sensitive data. | 922 | **N/A** |
| 8.2.3 | Verify that authenticated data is cleared from client storage, such as the  browser DOM, after the client or session is terminated. | 922 | **N/A** |

### V8.3 Sensitive Private Data

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 8.3.1 | Verify that sensitive data is sent to the server in the HTTP message body or  headers, and that query string parameters from any HTTP verb do not contain  sensitive data. | 319 | **Да** |
| 8.3.2 | Verify that users have a method to remove or export their data on demand.  | 212 | **Не** |
| 8.3.3 | Verify that users are provided clear language regarding collection and use of  supplied personal information and that users have provided opt-in consent  for the use of that data before it is used in any way. | 285 | **Не** |
| 8.3.4 | Verify that all sensitive data created and processed by the application has  been identified, and ensure that a policy is in place on how to deal with  sensitive data. | 200 | **Не** |
| 8.3.5 | Verify accessing sensitive data is audited (without logging the sensitive data  itself), if the data is collected under relevant data protection directives or  where logging of access is required. | 532 | **Да** |
| 8.3.6 | Verify that sensitive information contained in memory is overwritten as soon  as it is no longer required to mitigate memory dumping attacks, using zeroes  or random data. | 226 | **N/A** |
| 8.3.7 | Verify that sensitive or private information that is required to be encrypted, is  encrypted using approved algorithms that provide both confidentiality and  integrity.  | 327 | **Да** |
| 8.3.8 | Verify that sensitive personal information is subject to data retention  classification, such that old or out of date data is deleted automatically, on a  schedule, or as the situation requires. | 285 | **N/A** |

## V9 Communication

### V9.1 Client Communication Security

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 9.1.1 | Verify that TLS is used for all client connectivity, and does not fall back to  insecure or unencrypted communications. | 319 | **Не** |
| 9.1.2 | Verify using up to date TLS testing tools that only strong cipher suites are  enabled, with the strongest cipher suites set as preferred.| 326 | **Не** |
| 9.1.3 | Verify that only the latest recommended versions of the TLS protocol are  enabled, such as TLS 1.2 and TLS 1.3. The latest version of the TLS protocol should be the preferred option.| 326 | **Не** |

### V9.2 Server Communication Security

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 9.2.1 | Verify that connections to and from the server use trusted TLS certificates.  Where internally generated or self-signed certificates are used, the server  must be configured to only trust specific internal CAs and specific self-signed  certificates. All others should be rejected.| 295 | **Не** |
| 9.2.2 | Verify that encrypted communications such as TLS is used for all inbound and  outbound connections, including for management ports, monitoring,  authentication, API, or web service calls, database, cloud, serverless,  mainframe, external, and partner connections. The server must not fall back  to insecure or unencrypted protocols.| 319 | **Не** |
| 9.2.3 | Verify that all encrypted connections to external systems that involve  sensitive information or functions are authenticated.| 287 | **Не** |
| 9.2.4 | Verify that proper certification revocation, such as Online Certificate Status  Protocol (OCSP) Stapling, is enabled and configured.| 299 | **Не** |
| 9.2.5 | Verify that backend TLS connection failures are logged.| 544 | **Не** |

## V10 Malicious Code

### V10.1 Code Integrity

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 10.1.1 | Verify that a code analysis tool is in use that can detect potentially malicious  code, such as time functions, unsafe file operations and network  connections.| 749 | **Да**, коришћени су алати за анализу кода. |

### V10.2 Code Integrity

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 10.2.1 | Verify that the application source code and third party libraries do not  contain unauthorized phone home or data collection capabilities. Where  such functionality exists, obtain the user's permission for it to operate before  collecting any data.| 359 | **Да** |
| 10.2.2 | Verify that the application does not ask for unnecessary or excessive  permissions to privacy related features or sensors, such as contacts,  cameras, microphones, or location.| 272 | **Да** |
| 10.2.3 | Verify that the application source code and third party libraries do not  contain back doors, such as hard-coded or additional undocumented  accounts or keys, code obfuscation, undocumented binary blobs, rootkits, or  anti-debugging, insecure debugging features, or otherwise out of date,  insecure, or hidden functionality that could be used maliciously if  discovered.| 507 | **Вероватно**, апликација сигурно нема, али за библиотеке никад не можемо бити сигурни (пример: xz) |
| 10.2.4 | Verify that the application source code and third party libraries do not  contain time bombs by searching for date and time related functions.| 511 | **Вероватно** |
| 10.2.5 | Verify that the application source code and third party libraries do not  contain malicious code, such as salami attacks, logic bypasses, or logic  bombs.| 511 | **Вероватно** |
| 10.2.6 | Verify that the application source code and third party libraries do not  contain Easter eggs or any other potentially unwanted functionality.| 507 | **Вероватно** |

### V10.3 Code Integrity

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 10.3.1 | Verify that if the application has a client or server auto-update feature,  updates should be obtained over secure channels and digitally signed. The  update code must validate the digital signature of the update before  installing or executing the update.|  16 | **Да**, нема auto update |
| 10.3.2 | Verify that the application employs integrity protections, such as code  signing or subresource integrity. The application must not load or execute  code from untrusted sources, such as loading includes, modules, plugins,  code, or libraries from untrusted sources or the Internet.| 353 | **Не**  |
| 10.3.3 | Verify that the application has protection from subdomain takeovers if the  application relies upon DNS entries or DNS subdomains, such as expired  domain names, out of date DNS pointers or CNAMEs, expired projects at  public source code repos, or transient cloud APIs, serverless functions, or  storage buckets (autogen-bucket-id.cloud.example.com) or similar.  Protections can include ensuring that DNS names used by applications are  regularly checked for expiry or change.| 350 | **Да** |

## V11 Business Logic

### V11.1 Business Logic Security

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 11.1.1 | Verify that the application will only process business logic flows for the same  user in sequential step order and without skipping steps.| 841 | **Да** |
| 11.1.2 | Verify that the application will only process business logic flows with all steps  being processed in realistic human time, i.e. transactions are not submitted  too quickly.| 799 | **Не** |
| 11.1.3 | Verify the application has appropriate limits for specific business actions or  transactions which are correctly enforced on a per user basis.| 770 | **Не** |
| 11.1.4 | Verify that the application has anti-automation controls to protect against  excessive calls such as mass data exfiltration, business logic requests, file  uploads or denial of service attacks.| 770 | **Не** |
| 11.1.5 | Verify the application has business logic limits or validation to protect  against likely business risks or threats, identified using threat modeling or  similar methodologies.| 841 | **Не** |
| 11.1.6 | Verify that the application does not suffer from "Time Of Check to Time Of  Use" (TOCTOU) issues or other race conditions for sensitive operations.| 367 | **?** |
| 11.1.7 | Verify that the application monitors for unusual events or activity from a  business logic perspective. For example, attempts to perform actions out of  order or actions which a normal user would never attempt.| 754 | **Не** |
| 11.1.8 | Verify that the application has configurable alerting when automated attacks  or unusual activity is detected.| 390 | **Не** |

## V12 Files and Resources

### V12.1 File Upload

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 12.1.1 | Verify that the application will not accept large files that could fill up storage  or cause a denial of service.| 400 | **Не** |
| 12.1.2 | Verify that the application checks compressed files (e.g. zip, gz, docx, odt)  against maximum allowed uncompressed size and against maximum number  of files before uncompressing the file.| 409 | **Не**, не користе се компресовани фајлови |
| 12.1.3 | Verify that a file size quota and maximum number of files per user is  enforced to ensure that a single user cannot fill up the storage with too  many files, or excessively large files.| 770 | **Не** |

### V12.2 File Integrity

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 12.2.1 | Verify that files obtained from untrusted sources are validated to be of  expected type based on the file's content.| 434 | **Да** |

### V12.3 File Execution

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 12.3.1 | Verify that user-submitted filename metadata is not used directly by system  or framework filesystems and that a URL API is used to protect against path  traversal.|  22 | **Да**, не користимо овакве механизме |
| 12.3.2 | Verify that user-submitted filename metadata is validated or ignored to  prevent the disclosure, creation, updating or removal of local files (LFI).|  73 | **Даa** |
| 12.3.3 | Verify that user-submitted filename metadata is validated or ignored to  prevent the disclosure or execution of remote files via Remote File Inclusion  (RFI) or Server-side Request Forgery (SSRF) attacks.|  98 | **Да** |
| 12.3.4 | Verify that the application protects against Reflective File Download (RFD) by  validating or ignoring user-submitted filenames in a JSON, JSONP, or URL  parameter, the response Content-Type header should be set to text/plain,  and the Content-Disposition header should have a fixed filename.| 641 | **Не** |
| 12.3.5 | Verify that untrusted file metadata is not used directly with system API or  libraries, to protect against OS command injection.|  78 | **Да** |
| 12.3.6 | Verify that the application does not include and execute functionality from  untrusted sources, such as unverified content distribution networks,  JavaScript libraries, node npm libraries, or server-side DLLs.| 829 | **Да** |

### V12.4 File Storage

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 12.4.1 | Verify that files obtained from untrusted sources are stored outside the web  root, with limited permissions.| 552 | **Не** |
| 12.4.2 | Verify that files obtained from untrusted sources are scanned by antivirus  scanners to prevent upload and serving of known malicious content.| 509 | **Не** |

### V12.5 File Download

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 12.5.1 | Verify that the web tier is configured to serve only files with specific file  extensions to prevent unintentional information and source code leakage.  For example, backup files (e.g. .bak), temporary working files (e.g. .swp),  compressed files (.zip, .tar.gz, etc) and other extensions commonly used by  editors should be blocked unless required.| 552 | **N/A** |
| 12.5.2 | Verify that direct requests to uploaded files will never be executed as  HTML/JavaScript content.| 434 | **N/A** |

### V12.6 SSRF Protection

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 12.6.1 | Verify that the web or application server is configured with an allow list of  resources or systems to which the server can send requests or load data/files  from.| 918 | **a** |

## V13 API and Web Service

### V13.1 Generic Web Service Security

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 13.1.1 | Verify that all application components use the same encodings and parsers  to avoid parsing attacks that exploit different URI or file parsing behavior  that could be used in SSRF and RFI attacks.| 116 | **?** |
| 13.1.2 |  [DELETED, DUPLICATE OF 4.3.1] | | |
| 13.1.3 | Verify API URLs do not expose sensitive information, such as the API key,  session tokens etc.| 598 | **Да** |
| 13.1.4 | Verify that authorization decisions are made at both the URI, enforced by  programmatic or declarative security at the controller or router, and at the  resource level, enforced by model-based permissions.| 285 | **?** |
| 13.1.5 | Verify that requests containing unexpected or missing content types are  rejected with appropriate headers (HTTP response status 406 Unacceptable  or 415 Unsupported Media Type).| 434 | **Не** |

### V13.2 RESTful Web Service

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 13.2.1 | Verify that enabled RESTful HTTP methods are a valid choice for the user or  action, such as preventing normal users using DELETE or PUT on protected  API or resources.| 650 | **Да** |
| 13.2.2 | Verify that JSON schema validation is in place and verified before accepting  input.| 20  | **Да** |
| 13.2.3 | Verify that RESTful web services that utilize cookies are protected from  Cross-Site Request Forgery via the use of at least one or more of the  following: double submit cookie pattern, CSRF nonces, or Origin request  header checks.| 352 | **Да**, не користе се кукији |
| 13.2.4 | [DELETED, DUPLICATE OF 11.1.4]| | |
| 13.2.5 | Verify that REST services explicitly check the incoming Content-Type to be  the expected one, such as application/xml or application/json.| 436 | **Не** |
| 13.2.6 | Verify that the message headers and payload are trustworthy and not  modified in transit. Requiring strong encryption for transport (TLS only) may  be sufficient in many cases as it provides both confidentiality and integrity  protection. Per-message digital signatures can provide additional assurance  on top of the transport protections for high-security applications but bring  with them additional complexity and risks to weigh against the benefits.| 345 | **Не** |

### V13.3 SOAP Web Service

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 13.3.1 | Verify that XSD schema validation takes place to ensure a properly formed  XML document, followed by validation of each input field before any  processing of that data takes place.|  20 | **N/A** |
| 13.3.2 | Verify that the message payload is signed using WS-Security to ensure  reliable transport between client and service.| 345 | **Не** |

### V13.4 GraphQL

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 13.4.1 | Verify that a query allow list or a combination of depth limiting and amount  limiting is used to prevent GraphQL or data layer expression Denial of  Service (DoS) as a result of expensive, nested queries. For more advanced  scenarios, query cost analysis should be used.| 770 | **N/A** |
| 13.4.2 | Verify that GraphQL or other data layer authorization logic should be  implemented at the business logic layer instead of the GraphQL layer.| 285 | **N/A** |

## V14 Configuration

### V14.1 Build and Deploy

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 14.1.1 | Verify that the application build and deployment processes are performed in  a secure and repeatable way, such as CI / CD automation, automated  configuration management, and automated deployment scripts.|  | **Не** |
| 14.1.2 | Verify that compiler flags are configured to enable all available buffer  overflow protections and warnings, including stack randomization, data  execution prevention, and to break the build if an unsafe pointer, memory,  format string, integer, or string operations are found.| 120  | **Да**, Golang подразумевано нуди подршку |
| 14.1.3 | Verify that server configuration is hardened as per the recommendations of  the application server and frameworks in use.| 16 | **Не** |
| 14.1.4 | Verify that the application, configuration, and all dependencies can be re-deployed using automated deployment scripts, built from a documented and  tested runbook in a reasonable time, or restored from backups in a timely  fashion.|  | **Не** |
| 14.1.5 | Verify that authorized administrators can verify the integrity of all security-relevant configurations to detect tampering.|  | **Не** |

### V14.2 Build and Deploy

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 14.2.1 | Verify that all components are up to date, preferably using a dependency  checker during build or compile time.| 1026 | **Да**, али мануелно се ради |
| 14.2.2 | Verify that all unneeded features, documentation, sample applications and  configurations are removed.| 1002  | **Да**, мада sample апликације су део пројекта и сачуване су. |
| 14.2.3 | Verify that if application assets, such as JavaScript libraries, CSS or web  fonts, are hosted externally on a Content Delivery Network (CDN) or  external provider, Subresource Integrity (SRI) is used to validate the integrity  of the asset.| 829 | **Не** |
| 14.2.4 | Verify that third party components come from pre-defined, trusted and  continually maintained repositories.| 829 | **Да** |
| 14.2.5 | Verify that a Software Bill of Materials (SBOM) is maintained of all third  party libraries in use.|  | **Не** |
| 14.2.6 | Verify that the attack surface is reduced by sandboxing or encapsulating  third party libraries to expose only the required behaviour into the  application.| 265 | **Не** |

### V14.3 Unintended Security Disclosure

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 14.3.1 | [DELETED, DUPLICATE OF 7.4.1]|  |  |
| 14.3.2 | Verify that web or application server and application framework debug  modes are disabled in production to eliminate debug features, developer  consoles, and unintended security disclosures.| 497  | **N/A**, пројекат није постављан ван debug режима |
| 14.3.3 | Verify that the HTTP headers or any part of the HTTP response do not expose  detailed version information of system components.| 200 | **Да** |

### V14.4 HTTP Security Headers

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 14.4.1 | Verify that every HTTP response contains a Content-Type header. Also  specify a safe character set (e.g., UTF-8, ISO-8859-1) if the content types are  text/*, /+xml and application/xml. Content must match with the provided  Content-Type header.| 173 | **Да**, библиотеке (flask, go) раде ово уместо нас |
| 14.4.2 | Verify that all API responses contain a Content-Disposition: attachment;  filename="api.json" header (or other appropriate filename for the content  type).| 116  | **Не** |
| 14.4.3 | Verify that a Content Security Policy (CSP) response header is in place that  helps mitigate impact for XSS attacks like HTML, DOM, JSON, and JavaScript  injection vulnerabilities.| 1021 | **?** |
| 14.4.4 | Verify that all responses contain a X-Content-Type-Options: nosniff header. | 116 | **Не** |
| 14.4.5 | Verify that a Strict-Transport-Security header is included on all responses  and for all subdomains, such as Strict-Transport-Security: max-age=15724800; includeSubdomains.| 523  | **Не** |
| 14.4.6 | Verify that a suitable Referrer-Policy header is included to avoid exposing  sensitive information in the URL through the Referer header to untrusted  parties.| 116 | **Не** |
| 14.4.7 | Verify that the content of a web application cannot be embedded in a third-party site by default and that embedding of the exact resources is only  allowed where necessary by using suitable Content-Security-Policy: frame-ancestors and X-Frame-Options response head| 1021 | **?** |

### V14.5 HTTP Request Header Validation

| Број  | Опис | CWE | Испуњеност |
| ----- | ---- | --- | ---------- |
| 14.5.1 | Verify that the application server only accepts the HTTP methods in use by  the application/API, including pre-flight OPTIONS, and logs/alerts on any  requests that are not valid for the application context.| 749 | **Да** (парцијално) |
| 14.5.2 | Verify that the supplied Origin header is not used for authentication or  access control decisions, as the Origin header can easily be changed by an  attacker.| 346  | **Да** |
| 14.5.3 | Verify that the Cross-Origin Resource Sharing (CORS) Access-Control-Allow-Origin header uses a strict allow list of trusted domains and subdomains to  match against and does not support the "null" origin.| 346 | **Не** |
| 14.5.4 | Verify that HTTP headers added by a trusted proxy or SSO devices, such as a  bearer token, are authenticated by the application.| 306 | **N/A** |