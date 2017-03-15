[![Build Status](https://travis-ci.org/incu6us/asterisk-ami-api.svg)](https://travis-ci.org/incu6us/asterisk-ami-api)

## asterisk-ami-api

##### example:
```
curl -i "http://localhost:3000/api/v1/call/{SIPID}/{MSISDN}?async={async}"
```
- SIPID - sip internal number
- MSISDN - msisdn number
- async - is the call asynchronous ( `true` | `false` )