package crypto

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/models"
	"math/big"
)

const CollegeParkElectionKey = "FE211456D2A9E67D28009C885B9052B4999F2F97A930C2557AE9DD346BDD20F8FEB09FD8CD693D7FAC646A7683D028B43A87C224E62DC76832272223582C768CEDA25A17A6ADB0103CCD3B18175E0B226D54DA939C651497E091D1FEAB34EC45490E287E8A14C06397CEDD0FFA3C2B410ACD443E1961FCE0135ED1231499242B2E2A53A9618BB7C82BF76A0C035281A1C6E7ABCEEAADFDDC53969F350057C8723B96D00B3F9C72FF07AAE0BF94ED52D160812BDE79F131BD6708B30C2934E69E8C085F688B0692E5DCD908F984845DCDEAA33355C1B5D0283FFAF727D828D5D2DC7500A9342C3E0B85B7CD78BFA21B7D31A8FB4EF835DBD42BDFB9B423BDA11FCB552BC0B1FBD62CFD91444FB176FB9DB9B306A1ECC7944DD2E2CC269DA3C0A87E7CDB864838B5C385B04D6D6413505495F46A85C8C8AB04271186565C671B4050AF14871E4B996CC3261806F30C34C40222956B17830B299A1C988203BEF047B9585072281075F850B591BDB345FCE8C86FCFA68C2B5AFDFD4D170DB867466D722139A13D017DD783BC66050952951DC9002CB53C8ED5979E0794A6879076061A86E4248C7EF3AD91C50F4CA17ABB52E7AF6A376A34AB4AF9BF685B15035C8966F8B903FCCAB1B9D5BC5AB4AF0BE3EC4DDC822C7B683E61BBD1A84275398283356377998F2C7EF5A53C6133EF0A2C138AB0F9C1052C50F565F471C118F234FE"

// taken from: https://github.com/AU-HC/elliptic-curve-benchmark-go/blob/master/elgamal/elgamal.go

type PublicKey struct {
	G, P, K, Q *models.BigInt
}

type SecretKey struct {
	G, P, S *big.Int
}

func GenerateKeyPair() (PublicKey, SecretKey) {
	// Fixed G and P to ElectionGuard standard
	var publicKey PublicKey
	var secretKey SecretKey

	g := constants.GetG()
	p := constants.GetP()
	q := constants.GetQ()
	// s := GenerateRandomModQ()

	// secretKey.G = g
	// secretKey.P = p
	// secretKey.S = s

	publicKey.G = g
	publicKey.P = p
	publicKey.Q = q
	publicKey.K = models.MakeBigIntFromString(CollegeParkElectionKey, 16)

	return publicKey, secretKey
}

func Encrypt(publicKey PublicKey, m int, epsilon *models.BigInt) (models.BigInt, models.BigInt) {
	var alpha models.BigInt
	alpha.Exp(&publicKey.G.Int, &epsilon.Int, &publicKey.P.Int) // g^\epsilon mod p

	var kM big.Int
	var kR big.Int
	kM.Exp(&publicKey.K.Int, big.NewInt(int64(m)), &publicKey.P.Int) // K^m mod p
	kR.Exp(&publicKey.K.Int, &epsilon.Int, &publicKey.P.Int)         // K^\epsilon mod p

	var beta models.BigInt
	beta.Mul(&kM, &kR)                    // K^m K^\epsilon
	beta.Mod(&beta.Int, &publicKey.P.Int) // K^m \cdot K^\epsilon mod p

	return alpha, beta // (g^\epsilon, K^m \cdot K^\epsilon mod p)
}
