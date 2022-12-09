# BitBucket to Telegram PR notifications bridge

This is a simple Telegram bot that delivers BitBucket PR events to users.
Bridge might be useful if a BitBucket instance is hidden behind a VPN, and it is hard for a team to track PR updates.

IMPORTANT:
* Only the BitBucket API v1 is supported;
* Suitable for small teams with low PR activity since the implementation is pretty basic;
* For now the bot only support one BitBucket project;

## Configuration

Refer to [config pkg](./pkg/config) to see the list of available Viper keys.
Every key can be set via environment variable using the form:

```
BBTT__<SECTION>__<KEY>
```
