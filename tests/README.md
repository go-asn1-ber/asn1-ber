# Conformance test fixtures

The `tc*.ber` files are BER-encoded inputs from the Strozhevsky free ASN.1 test
suite, used by `TestSuiteDecodePacket` / `TestSuiteReadPacket` in
`../suite_test.go` to exercise the decoder against well-formed and deliberately
malformed encodings.

- Descriptions: <http://www.strozhevsky.com/free_docs/free_asn1_testsuite_descr.pdf>
- Original archive: <http://www.strozhevsky.com/free_docs/TEST_SUITE.zip>

Each fixture's expected decode outcome (success, expected error string, or
abnormal/indefinite re-encoding) is recorded in the `testCases` table in
`../suite_test.go`.
