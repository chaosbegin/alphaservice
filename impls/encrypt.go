package impls

import "github.com/Luzifer/go-openssl/v3"

const pwdKey = "pwdKey1"
const pwdKey2 = "pwdKey2"

func PwdDecrypt(pwd string) (string, error) {
	o := openssl.New()
	out, err := o.DecryptBytes(pwdKey, []byte(pwd), openssl.DigestMD5Sum)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func PwdEncrypt(pwd string) (string, error) {
	o := openssl.New()
	out, err := o.EncryptBytes(pwdKey, []byte(pwd), openssl.DigestMD5Sum)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func PwdDecrypt2(pwd string) (string, error) {
	o := openssl.New()
	out, err := o.DecryptBytes(pwdKey2, []byte(pwd), openssl.DigestSHA256Sum)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func PwdEncrypt2(pwd string) (string, error) {
	o := openssl.New()
	out, err := o.EncryptBytes(pwdKey2, []byte(pwd), openssl.DigestSHA256Sum)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
