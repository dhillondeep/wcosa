package semver

func (a *Version) Eq(b *Version) bool {
    return a.eq(b)
}

func (a *Version) Ne(b *Version) bool {
    return !a.eq(b)
}

func (a *Version) Lt(b *Version) bool {
    return a.less(b)
}

func (a *Version) Gt(b *Version) bool {
    return b.less(a)
}

func (a *Version) Le(b *Version) bool {
    return !b.less(a)
}

func (a *Version) Ge(b *Version) bool {
    return !a.less(b)
}
