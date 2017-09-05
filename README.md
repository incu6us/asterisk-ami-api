## asterisk-ami-api [![Build Status](https://travis-ci.org/incu6us/asterisk-ami-api.svg)](https://travis-ci.org/incu6us/asterisk-ami-api)

##### example:
- call:
    ```
    curl -i "http://localhost:3000/api/v1/call/{SIPID}/{MSISDN}?async={async}"
    ```
    - SIPID - sip internal number
    - MSISDN - msisdn number
    - async - asynchronous call. default: `false` ( value: `true` || `false` )

- playback:
    ```
    curl -i "http://localhost:3000/api/v1/playback/{MSISDN}/{FILE-TO-PLAY}?async={async}"
    ```
    - MSISDN - msisdn number
    - FILE-TO-PLAY - file in gsm-format
    - async - asynchronous call. default: `false` ( value: `true` || `false` )


- sms:
    ```
    curl -XPOST "http://localhost:3000/api/v1/modem/send/sms/{modem}/{MSISDN}" --data "test message"
    ```
    
- cdr:
    ```
    curl -XGET "http://localhost:3000/api/v1/cdr/search/0937530214"
    ```
    or
    ```
    curl -XGET "http://localhost:3000/api/v1/cdr/search/0937530214?startdate=2017-04-26&enddate=2017-04-26"
    ```

#### External vars(for [Dockerized image](https://hub.docker.com/r/incu6us/asterisk-ami-api)):
##### AMI config
- *AMI_HOST* - host, where AMI interface is located
- *AMI_PORT* - AMI port
- *AMI_USER* - username
- *AMI_PASS* - password

##### database config
- *DB_HOST* - database host ip
- *DB_DBNAME* - database name
- *DB_USER* - username
- *DB_PASS* - password

##### asterisk dialplan config
- *ASTERISK_CONTEXT* - main dial context
- *ASTERISK_PLAYBACK_CONTEXT* - dialplan with playback