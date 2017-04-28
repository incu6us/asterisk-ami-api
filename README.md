[![Build Status](https://travis-ci.org/incu6us/asterisk-ami-api.svg)](https://travis-ci.org/incu6us/asterisk-ami-api)

## asterisk-ami-api

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
    curl -XPOST "http://localhost:3000/api/v1/cdr/search/0937530214"
    ```
    or
    ```
    curl -XPOST "http://localhost:3000/api/v1/cdr/search/0937530214?startdate=2017-04-26&enddate=2017-04-26"
    ```