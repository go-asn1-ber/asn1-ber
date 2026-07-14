# Security Policy

This library decodes ASN.1 BER, frequently from untrusted network input (for
example, LDAP messages). Parsing bugs can have security impact, so reports are
welcome.

## Reporting a vulnerability

Please report suspected vulnerabilities privately using GitHub's
[private vulnerability reporting](https://github.com/go-asn1-ber/asn1-ber/security/advisories/new)
rather than opening a public issue. Include a description, affected versions, and
a reproducing input (a BER byte sequence) where possible.

## Hardening notes for callers

When decoding untrusted data, keep the default limits in place or set your own:

- `MaxPacketLengthBytes` bounds the size of any single decoded packet.
- `MaxNestingDepth` bounds recursion into constructed packets.

Setting either to `0` disables that limit and is not recommended for untrusted
input.
