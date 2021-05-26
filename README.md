# genshindaily

![Discord Screenshot](https://i.imgur.com/M3jVlR0.png)

Yet another Genshin auto daily check-in tool.
Based on [genshinhelper by y1ndan](https://github.com/y1ndan/genshinhelper),
but for HoYoLab check-ins only.

## Running

### AWS

You will need to set up a Lambda function and then trigger it with a Cloud Events cronjob.
The cronjob should also be set up to push the configuration as a constant JSON string.

### Local

You can build a runnable version with `make build-cmd` and then execute it by piping the configuration through stdin, e.g. `./daily < conf.json`

## Example Configuration

```json
{
    "cookies": "mi18nLang=en-us; account_id=YOUR_ACCOUNT_ID_HERE; cookie_token=YOUR_TOKEN_HERE",
    "notifiers": [
        {
            "kind": "discord",
            "url": "https://discord.com/api/webhooks/123456789012345678/bar"
        },
        {
            "kind": "health_check",
            "url": "https://hc-ping.com/foo"
        }
    ]
}
```

## Notifications

Currently supports Discord webhooks and healthchecks.io notifications only.
