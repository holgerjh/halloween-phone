# Example barsip configuration

**WARNING: Make sure to double-check the contact that is dialed if you run
`/dialcontact` so you do not perform calls that may cost you money. If
possible, block outgoing calls in your SIP provider. I do not take any
responsibility for costs that occure using this program. Use at your own risk!**

Place these files into your ~/.baresip folder (create if needed) and change
them as follows:

* `accounts`: Create a SIP phone in your sip provider (e.g. FritzBox) with name
  `halloween-phone`. Change `SIP_IP` to the ip of your sip provider, replace
  `PASSWORD` with the account's password.

* `config`: Adapt as needed. Make sure the paths are ok and the source
  encodings matches. 

* `contacts`: Replace `NUMBER` and `SIP_IP`, e.g. with `**1` (local connection)
  and `192.168.1.1`.

You should then start baresip and do a `/dialcontact`. If it does not work
perform a `/dialcontact target`. After that, check that `/dialcontact` works. 
