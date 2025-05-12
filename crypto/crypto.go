package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"math/big"
)

type CurvePoint struct {
	X *big.Int
	Y *big.Int
}

type CryptoContext struct {
	Curve elliptic.Curve
	H     CurvePoint
	G     CurvePoint
}

func NewContext() *CryptoContext {
	var context CryptoContext

	context.Curve = elliptic.P256()
	context.G = GenerateRandomPoint(context.Curve)

	return &context
}

func (context *CryptoContext) GeneratePublicKey(secretNum *big.Int) CurvePoint {
	curve := elliptic.P256()

	priv, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Fatal(err)
	}

	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarMult(context.G.X, context.G.Y, secretNum.Bytes())

	return CurvePoint{priv.PublicKey.X, priv.PublicKey.Y}
}

func (context *CryptoContext) GenerateAggregateKey(publicKeys []CurvePoint) {
	curve := elliptic.P256()

	aggregateX, aggregateY := curve.ScalarBaseMult(big.NewInt(0).Bytes())

	for _, pub := range publicKeys {
		aggregateX, aggregateY = curve.Add(aggregateX, aggregateY, pub.X, pub.Y)
	}

	context.H = CurvePoint{aggregateX, aggregateY}
}

func GenerateRandomPoint(curve elliptic.Curve) CurvePoint {
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	x, y := priv.PublicKey.X, priv.PublicKey.Y

	return CurvePoint{
		x,
		y,
	}
}

func (context *CryptoContext) GenerateRandomPoint() CurvePoint {
	priv, err := ecdsa.GenerateKey(context.Curve, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	x, y := priv.PublicKey.X, priv.PublicKey.Y

	return CurvePoint{
		x,
		y,
	}
}

func (context *CryptoContext) MaskCard(M CurvePoint, r *big.Int) (CurvePoint, CurvePoint) {
	Cax, Cay := context.Curve.ScalarMult(context.G.X, context.G.Y, r.Bytes())
	rHx, rHy := context.Curve.ScalarMult(context.H.X, context.H.Y, r.Bytes())
	Cbx, Cby := context.Curve.Add(M.X, M.Y, rHx, rHy)

	return CurvePoint{Cax, Cay}, CurvePoint{Cbx, Cby}
}

func (context *CryptoContext) RemaskCard(Ca CurvePoint, Cb CurvePoint, r *big.Int) (CurvePoint, CurvePoint) {
	rpGx, rpGy := context.Curve.ScalarMult(context.G.X, context.G.Y, r.Bytes())
	rpHx, rpHy := context.Curve.ScalarMult(context.H.X, context.H.Y, r.Bytes())
	Cpax, Cpay := context.Curve.Add(Ca.X, Ca.Y, rpGx, rpGy)
	Cpbx, Cpby := context.Curve.Add(Cb.X, Cb.Y, rpHx, rpHy)

	return CurvePoint{Cpax, Cpay}, CurvePoint{Cpbx, Cpby}
}

func (context *CryptoContext) DecryptCard(Ca CurvePoint, Cb CurvePoint, xi *big.Int) CurvePoint {
	Cax, Cay := context.Curve.ScalarMult(Ca.X, Ca.Y, xi.Bytes())
	negCay := new(big.Int).Neg(Cay)
	negCay = negCay.Mod(negCay, context.Curve.Params().P)

	Cbx, Cby := context.Curve.Add(Cb.X, Cb.Y, Cax, negCay)

	return CurvePoint{Cbx, Cby}
}
