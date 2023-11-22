package game

type DamageCounter struct {
	percentage  float32
	damageTaken float32
	duration    float32
}

func NewDamageCounter(duration float32) DamageCounter {
	return DamageCounter{
		percentage:  0.,
		damageTaken: 0.,
		duration:    duration,
	}
}

func DamageCalc(damage DamageCounter, damageStep float32) DamageCounter {
	damage_ := damage
	damage_.damageTaken += damageStep
	damage_.percentage = (damage.damageTaken / damage.duration) * 100.

	if damage_.percentage >= 100. {
		damage_.percentage = 100.
	}

	return damage_
}

func MaxDamage(damage DamageCounter) bool {
	return damage.percentage == 100.
}
