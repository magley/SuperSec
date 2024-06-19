У овом документу налази се анализа OWASP ASVS 4.0.3 документа примењен на Lesotho пројекат.

## V1 Architecture, Design and Threat Modeling

### V1.1 Secure Software Development Lifecycle

Није разматрано. Овај пројекат је сувише малих размера и разбијање на целине,
корисничке приче итд.; а ради дефинисања безбедносних апсеката није исплативо.

### V1.2 Authentication Architecture

| Број | Опис | CWE | Испуњеност |
| ---- | ---- | --- | ---------- |
| 1.2.1| Verify the use of unique or special low-privilege operating system accounts for all application components, services, and servers.| 250 | **Да**, ниједан сервис не захтева администраторске привилегије. |
| 1.2.2| Verify that communications between application components, including APIs, middleware and data layers, are authenticated. Components should have the least necessary privileges needed. | 306 | **Не**, интерне компоненте су подешене са подразумеваним привилегијама. Свако може приступити Lesotho севрису - **није имплементиран механизам API кључева**. |
| 1.2.3| Verify that the application uses a single vetted authentication mechanism that is known to be secure, can be extended to include strong authentication, and has sufficient logging and monitoring to detect account abuse or breaches. | 306 | **Не**, logging и monitoring нису имплементирани. |
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
| 1.7.1 | Verify that a common logging format and approach is used across the system. | 1009 | **Не**, логовање није имплементирано. |
| 1.7.2 | Verify that logs are securely transmitted to a preferably remote system for analysis, detection, alerting, and escalation. | | **Не**, логовање није имплементирано. |

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
| 4.1.1 | Verify that the application enforces access control rules on a trusted service layer, especially if client-side access control is present and could be bypassed. | 602 | ** ** |
| 4.1.2 | | | |
| 4.1.3 | | | |
| 4.1.4 | | | |
| 4.1.5 | | | | 