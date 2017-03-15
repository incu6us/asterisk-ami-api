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


- sms:
    ```
    curl -XPOST "http://localhost:3000/api/v1/modem/send/sms/{modem}/{MSISDN}" --data "test message"
    ```