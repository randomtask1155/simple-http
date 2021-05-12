

# root handler

responds with html response

# get/data?length=200

```
curl "https://simple-http.cfapps-48.slot-59.pez.vmware.com/get/data?length=1000000" -k
```

returns a json response with given character length 

```
{"data": "XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTeMaPEZQleQYhYzRyWJjPjzpfRFEgmotaFetHsbZRjxAwnwekrBEmfdzdcEkXBAkjQZLCtTMtTCoaNatyyiNKAReKJyiXJrscctNswYNsGRussVmaozFZBsbOJiFQGZsnwTKSmVoiG"}
```

# post/data

```
curl -X POST "https://simple-http.cfapps-48.slot-59.pez.vmware.com/post/data" -k -F 'file=@/Users/danl/Downloads/ubuntu-20.04-desktop-amd64.iso'
```

reads request body and exits