package main

import "sort"

func (ll *LineList) DoClip(fl *FaceList) {
	if ll.count < 1 {
		return
	}

	faceCount := ll.count
	dl := ll.dl

	// clip Z near
	for t := 0; t < faceCount; t++ {
		l := dl[t]
		if l.visible {
			if l.z1 < nearZ {
				if l.z2 < nearZ {
					l.visible = false
				} else {
					fakt := (nearZ - l.z1) / (l.z1 - l.z2)
					l.x1 += (l.x1 - l.x2) * fakt
					l.y1 += (l.y1 - l.y2) * fakt
					l.z1 = nearZ
				}
			} else {
				if l.z2 < nearZ {
					fakt := (nearZ - l.z2) / (l.z2 - l.z1)
					l.x2 += (l.x2 - l.x1) * fakt
					l.y2 += (l.y2 - l.y1) * fakt
					l.z2 = nearZ
				}
			}
		}
	}

	// clip Z far
	for t := 0; t < faceCount; t++ {
		l := dl[t]
		if l.visible {
			if l.z1 > farZ {
				if l.z2 > farZ {
					l.visible = false
				} else {
					fakt := (farZ - l.z1) / (l.z1 - l.z2)
					l.x1 += (l.x1 - l.x2) * fakt
					l.y1 += (l.y1 - l.y2) * fakt
					l.z1 = farZ
				}
			} else {
				if l.z2 > farZ {
					fakt := (farZ - l.z2) / (l.z2 - l.z1)
					l.x2 += (l.x2 - l.x1) * fakt
					l.y2 += (l.y2 - l.y1) * fakt
					l.z2 = farZ
				}
			}
		}
	}

	// clip X left
	for t := 0; t < faceCount; t++ {
		l := dl[t]
		if l.visible {
			if l.x1 < -l.z1 {
				if l.x2 < -l.z2 {
					l.visible = false
				} else {
					fakt := (l.z2 - l.z1)
					skal := (-l.z1 - l.x1) / (l.x2 - l.x1 + fakt)
					neu := skal * fakt
					l.y1 += (l.y2 - l.y1) * skal
					l.z1 += neu
					l.x1 = -l.z1
				}
			} else {
				if l.x2 < -l.z2 {
					fakt := (l.z1 - l.z2)
					skal := (-l.z2 - l.x2) / (l.x1 - l.x2 + fakt)
					neu := skal * fakt
					l.y2 += (l.y1 - l.y2) * skal
					l.z2 += neu
					l.x2 = -l.z2
				}
			}
		}
	}

	// clip X right
	for t := 0; t < faceCount; t++ {
		l := dl[t]
		if l.visible {
			if l.x1 > l.z1 {
				if l.x2 > l.z2 {
					l.visible = false
				} else {
					fakt := (l.z2 - l.z1)
					skal := (l.z1 - l.x1) / (l.x2 - l.x1 - fakt)
					neu := skal * fakt
					l.y1 += (l.y2 - l.y1) * skal
					l.z1 += neu
					l.x1 = l.z1
				}
			} else {
				if l.x2 > l.z2 {
					fakt := (l.z1 - l.z2)
					skal := (l.z2 - l.x2) / (l.x1 - l.x2 - fakt)
					neu := skal * fakt
					l.y2 += (l.y1 - l.y2) * skal
					l.z2 += neu
					l.x2 = l.z2
				}
			}
		}
	}

	// clip Y bottom
	for t := 0; t < faceCount; t++ {
		l := dl[t]
		if l.visible {
			if l.y1 < -l.z1 {
				if l.y2 < -l.z2 {
					l.visible = false
				} else {
					fakt := (l.z2 - l.z1)
					skal := (-l.z1 - l.y1) / (l.y2 - l.y1 + fakt)
					neu := fakt * skal
					l.x1 += (l.x2 - l.x1) * skal
					l.z1 += neu
					l.y1 = -l.z1
				}
			} else {
				if l.y2 < -l.z2 {
					fakt := (l.z1 - l.z2)
					skal := (-l.z2 - l.y2) / (l.y1 - l.y2 + fakt)
					neu := fakt * skal
					l.x2 += (l.x1 - l.x2) * skal
					l.z2 += neu
					l.y2 = -l.z2
				}
			}
		}
	}

	// clip Y top
	for t := 0; t < faceCount; t++ {
		l := dl[t]
		if l.visible {
			if l.y1 > l.z1 {
				if l.y2 > l.z2 {
					l.visible = false
				} else {
					fakt := (l.z2 - l.z1)
					skal := (l.z1 - l.y1) / (l.y2 - l.y1 - fakt)
					neu := fakt * skal
					l.x1 += (l.x2 - l.x1) * skal
					l.z1 += neu
					l.y1 = l.z1
				}
			} else {
				if l.y2 > l.z2 {
					fakt := (l.z1 - l.z2)
					skal := (l.z2 - l.y2) / (l.y1 - l.y2 - fakt)
					neu := fakt * skal
					l.x2 += (l.x1 - l.x2) * skal
					l.z2 += neu
					l.y2 = l.z2
				}
			}
		}
	}

	// remove not visible and project
	d_anf := 0
	for t := 0; t < ll.count; t++ {
		l := dl[t]
		if l.visible {
			if l.z1 > l.z2 {
				l.maxZ = l.z1
				l.minZ = l.z2
			} else {
				l.maxZ = l.z2
				l.minZ = l.z1
			}
			z1 := 1.0 / l.z1
			z2 := 1.0 / l.z2
			x1 := l.x1 * z1
			y1 := l.y1 * z1
			x2 := l.x2 * z2
			y2 := l.y2 * z2
			l.x1 = (320.0 + 319.0*x1)
			l.y1 = (240.0 + 239.0*y1)
			l.x2 = (320.0 + 319.0*x2)
			l.y2 = (240.0 + 239.0*y2)

			dl[d_anf] = l
			d_anf++
		}
	}
	ll.count = d_anf

	// occlusion
	clipperfect := false
	for !clipperfect {
		clipperfect = true
		for c := fl.count - 1; c >= 0; c-- {
			f := fl.dl[c]
			if !f.visible {
				continue
			}
			s := ll.count
			for t := 0; t < s; t++ {
				l := dl[t]
				if l.visible && l.maxZ > f.minZ {
					bo := true
					if (l.x1 < f.minx) && (l.x2 < f.minx) {
						bo = false
					}
					if (l.x1 > f.maxx) && (l.x2 > f.maxx) {
						bo = false
					}
					if (l.y1 < f.miny) && (l.y2 < f.miny) {
						bo = false
					}
					if (l.y1 > f.maxy) && (l.y2 > f.maxy) {
						bo = false
					}

					if bo {
						for i := 0; i < f.lcount && bo; i++ {
							bo = f.l[i] != l
						}
					}

					const eps = 1e-5
					const imeps = 1.0 - eps

					if bo {
						b0 := 1.0 / (((f.x2-f.x1)*(f.y3-f.y1)) - ((f.x3-f.x1)*(f.y2-f.y1)))
						b1 := (((f.x2-l.x1)*(f.y3-l.y1)) - ((f.x3-l.x1)*(f.y2-l.y1))) * b0
						b2 := (((f.x3-l.x1)*(f.y1-l.y1)) - ((f.x1-l.x1)*(f.y3-l.y1))) * b0
						b3 := (((f.x1-l.x1)*(f.y2-l.y1)) - ((f.x2-l.x1)*(f.y1-l.y1))) * b0
						ins1 := (b1 > -eps) && (b2 > -eps) && (b3 > -eps)

						b1 = (((f.x2-l.x2)*(f.y3-l.y2)) - ((f.x3-l.x2)*(f.y2-l.y2))) * b0
						b2 = (((f.x3-l.x2)*(f.y1-l.y2)) - ((f.x1-l.x2)*(f.y3-l.y2))) * b0
						b3 = (((f.x1-l.x2)*(f.y2-l.y2)) - ((f.x2-l.x2)*(f.y1-l.y2))) * b0
						ins2 := (b1 > -eps) && (b2 > -eps) && (b3 > -eps)

						y4y3 := l.y2 - l.y1
						x2x1 := f.x2 - f.x1
						y2y1 := f.y2 - f.y1
						x4x3 := l.x2 - l.x1
						q1 := (y4y3 * x2x1) - (x4x3 * y2y1)
						b1 = -1
						if (q1 != 0.0) && ((x2x1 != 0.0) || (y2y1 != 0.0)) {
							y3y1 := l.y1 - f.y1
							x3x1 := l.x1 - f.x1
							b1 = ((x3x1*y2y1)-(y3y1*x2x1)) / q1
							var a1 float64
							if x2x1 != 0 {
								a1 = (x3x1 + (b1 * x4x3)) / x2x1
							} else {
								a1 = (y3y1 + (b1 * y4y3)) / y2y1
							}
							if (a1 < eps) || (a1 > imeps) {
								b1 = -1
							}
						}

						y4y3 = l.y2 - l.y1
						x2x1 = f.x3 - f.x2
						y2y1 = f.y3 - f.y2
						x4x3 = l.x2 - l.x1
						q2 := (y4y3 * x2x1) - (x4x3 * y2y1)
						b2 = -1
						if (q2 != 0.0) && ((x2x1 != 0.0) || (y2y1 != 0.0)) {
							y3y1 := l.y1 - f.y2
							x3x1 := l.x1 - f.x2
							b2 = ((x3x1*y2y1)-(y3y1*x2x1)) / q2
							var a2 float64
							if x2x1 != 0 {
								a2 = (x3x1 + (b2 * x4x3)) / x2x1
							} else {
								a2 = (y3y1 + (b2 * y4y3)) / y2y1
							}
							if (a2 < eps) || (a2 > imeps) {
								b2 = -1
							}
						}

						y4y3 = l.y2 - l.y1
						x2x1 = f.x1 - f.x3
						y2y1 = f.y1 - f.y3
						x4x3 = l.x2 - l.x1
						q3 := (y4y3 * x2x1) - (x4x3 * y2y1)
						b3 = -1
						if (q3 != 0.0) && ((x2x1 != 0.0) || (y2y1 != 0.0)) {
							y3y1 := l.y1 - f.y3
							x3x1 := l.x1 - f.x3
							b3 = ((x3x1*y2y1)-(y3y1*x2x1)) / q3
							var a3 float64
							if x2x1 != 0 {
								a3 = (x3x1 + (b3 * x4x3)) / x2x1
							} else {
								a3 = (y3y1 + (b3 * y4y3)) / y2y1
							}
							if (a3 < eps) || (a3 > imeps) {
								b3 = -1
							}
						}

						bo1 := (b1 >= eps) && (b1 <= imeps)
						bo2 := (b2 >= eps) && (b2 <= imeps)
						bo3 := (b3 >= eps) && (b3 <= imeps)

						if ins1 {
							if ins2 {
								l.visible = false
							} else {
								if bo1 && !bo2 && !bo3 {
									l.x1 = l.x1 + b1*(l.x2-l.x1)
									l.y1 = l.y1 + b1*(l.y2-l.y1)
								} else if bo2 && !bo1 && !bo3 {
									l.x1 = l.x1 + b2*(l.x2-l.x1)
									l.y1 = l.y1 + b2*(l.y2-l.y1)
								} else if bo3 && !bo1 && !bo2 {
									l.x1 = l.x1 + b3*(l.x2-l.x1)
									l.y1 = l.y1 + b3*(l.y2-l.y1)
								}
							}
						} else {
							if ins2 {
								if bo1 && !bo2 && !bo3 {
									l.x2 = l.x1 + b1*(l.x2-l.x1)
									l.y2 = l.y1 + b1*(l.y2-l.y1)
								} else if bo2 && !bo1 && !bo3 {
									l.x2 = l.x1 + b2*(l.x2-l.x1)
									l.y2 = l.y1 + b2*(l.y2-l.y1)
								} else if bo3 && !bo1 && !bo2 {
									l.x2 = l.x1 + b3*(l.x2-l.x1)
									l.y2 = l.y1 + b3*(l.y2-l.y1)
								}
							} else {
								if bo1 && bo2 && !bo3 {
									if ll.scratchCount < len(ll.scratchLines) && ll.count < len(dl) {
										if b2 < b1 {
											b1, b2 = b2, b1
										}
										l1 := &ll.scratchLines[ll.scratchCount]
										ll.scratchCount++
										l1.minZ = l.minZ
										l1.maxZ = l.maxZ
										l1.visible = true
										l1.x2 = l.x2
										l1.y2 = l.y2
										l1.x1 = l.x1 + b2*(l.x2-l.x1)
										l1.y1 = l.y1 + b2*(l.y2-l.y1)
										l.x2 = l.x1 + b1*(l.x2-l.x1)
										l.y2 = l.y1 + b1*(l.y2-l.y1)
										dl[ll.count] = l1
										ll.count++
									}
								} else if bo1 && bo3 && !bo2 {
									if ll.scratchCount < len(ll.scratchLines) && ll.count < len(dl) {
										if b3 < b1 {
											b1, b3 = b3, b1
										}
										l1 := &ll.scratchLines[ll.scratchCount]
										ll.scratchCount++
										l1.minZ = l.minZ
										l1.maxZ = l.maxZ
										l1.visible = true
										l1.x2 = l.x2
										l1.y2 = l.y2
										l1.x1 = l.x1 + b3*(l.x2-l.x1)
										l1.y1 = l.y1 + b3*(l.y2-l.y1)
										l.x2 = l.x1 + b1*(l.x2-l.x1)
										l.y2 = l.y1 + b1*(l.y2-l.y1)
										dl[ll.count] = l1
										ll.count++
									}
								} else if bo2 && bo3 && !bo1 {
									if ll.scratchCount < len(ll.scratchLines) && ll.count < len(dl) {
										if b3 < b2 {
											b2, b3 = b3, b2
										}
										l1 := &ll.scratchLines[ll.scratchCount]
										ll.scratchCount++
										l1.minZ = l.minZ
										l1.maxZ = l.maxZ
										l1.visible = true
										l1.x2 = l.x2
										l1.y2 = l.y2
										l1.x1 = l.x1 + b3*(l.x2-l.x1)
										l1.y1 = l.y1 + b3*(l.y2-l.y1)
										l.x2 = l.x1 + b2*(l.x2-l.x1)
										l.y2 = l.y1 + b2*(l.y2-l.y1)
										dl[ll.count] = l1
										ll.count++
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (fl *FaceList) DoClip(mz bool) {
	if fl.count < 1 {
		return
	}

	dl := fl.dl
	// clip Z near
	faceCount := fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		f.visible = true
		if f.visible {
			if f.z1 < nearZ {
				if f.z2 < nearZ {
					if f.z3 < nearZ {
						f.visible = false
					} else {
						fakt := (nearZ - f.z1) / (f.z1 - f.z3)
						f.x1 += (f.x1 - f.x3) * fakt
						f.y1 += (f.y1 - f.y3) * fakt
						f.z1 = nearZ
						f.u1 += (f.u1 - f.u3) * fakt
						f.v1 += (f.v1 - f.v3) * fakt
						fakt = (nearZ - f.z2) / (f.z2 - f.z3)
						f.x2 += (f.x2 - f.x3) * fakt
						f.y2 += (f.y2 - f.y3) * fakt
						f.z2 = nearZ
						f.u2 += (f.u2 - f.u3) * fakt
						f.v2 += (f.v2 - f.v3) * fakt
					}
				} else {
					if f.z3 < nearZ {
						fakt := (nearZ - f.z1) / (f.z1 - f.z2)
						f.x1 += (f.x1 - f.x2) * fakt
						f.y1 += (f.y1 - f.y2) * fakt
						f.z1 = nearZ
						f.u1 += (f.u1 - f.u2) * fakt
						f.v1 += (f.v1 - f.v2) * fakt
						fakt = (nearZ - f.z3) / (f.z3 - f.z2)
						f.x3 += (f.x3 - f.x2) * fakt
						f.y3 += (f.y3 - f.y2) * fakt
						f.z3 = nearZ
						f.u3 += (f.u3 - f.u2) * fakt
						f.v3 += (f.v3 - f.v2) * fakt
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f1.z1 - nearZ) / (f1.z1 - f1.z2)
						f1.x1 += (f1.x2 - f1.x1) * fakt
						f1.y1 += (f1.y2 - f1.y1) * fakt
						f1.z1 = nearZ
						f1.u1 += (f1.u2 - f1.u1) * fakt
						f1.v1 += (f1.v2 - f1.v1) * fakt
						f1.x3 = f.x1
						f1.y3 = f.y1
						f1.z3 = nearZ
						f1.u3 = f.u1
						f1.v3 = f.v1
						fakt = (f.z1 - nearZ) / (f.z1 - f.z3)
						f.x1 += (f.x3 - f.x1) * fakt
						f.y1 += (f.y3 - f.y1) * fakt
						f.z1 = nearZ
						f.u1 += (f.u3 - f.u1) * fakt
						f.v1 += (f.v3 - f.v1) * fakt
					}
				}
			} else {
				if f.z2 < nearZ {
					if f.z3 < nearZ {
						fakt := (nearZ - f.z2) / (f.z2 - f.z1)
						f.x2 += (f.x2 - f.x1) * fakt
						f.y2 += (f.y2 - f.y1) * fakt
						f.z2 = nearZ
						f.u2 += (f.u2 - f.u1) * fakt
						f.v2 += (f.v2 - f.v1) * fakt
						fakt = (nearZ - f.z3) / (f.z3 - f.z1)
						f.x3 += (f.x3 - f.x1) * fakt
						f.y3 += (f.y3 - f.y1) * fakt
						f.z3 = nearZ
						f.u3 += (f.u3 - f.u1) * fakt
						f.v3 += (f.v3 - f.v1) * fakt
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - nearZ) / (f.z2 - f.z3)
						f.x2 += (f.x3 - f.x2) * fakt
						f.y2 += (f.y3 - f.y2) * fakt
						f.z2 = nearZ
						f.u2 += (f.u3 - f.u2) * fakt
						f.v2 += (f.v3 - f.v2) * fakt
						fakt = (f1.z2 - nearZ) / (f1.z2 - f1.z1)
						f1.x2 += (f1.x1 - f1.x2) * fakt
						f1.y2 += (f1.y1 - f1.y2) * fakt
						f1.z2 = nearZ
						f1.u2 += (f1.u1 - f1.u2) * fakt
						f1.v2 += (f1.v1 - f1.v2) * fakt
						f1.x3 = f.x2
						f1.y3 = f.y2
						f1.z3 = nearZ
						f1.u3 = f.u2
						f1.v3 = f.v2
					}
				} else {
					if f.z3 < nearZ {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z3 - nearZ) / (f.z3 - f.z2)
						f.x3 += (f.x2 - f.x3) * fakt
						f.y3 += (f.y2 - f.y3) * fakt
						f.z3 = nearZ
						f.u3 += (f.u2 - f.u3) * fakt
						f.v3 += (f.v2 - f.v3) * fakt
						fakt = (f1.z3 - nearZ) / (f1.z3 - f1.z1)
						f1.x3 += (f1.x1 - f1.x3) * fakt
						f1.y3 += (f1.y1 - f1.y3) * fakt
						f1.z3 = nearZ
						f1.u3 += (f1.u1 - f1.u3) * fakt
						f1.v3 += (f1.v1 - f1.v3) * fakt
						f1.x2 = f.x3
						f1.y2 = f.y3
						f1.z2 = nearZ
						f1.u2 = f.u3
						f1.v2 = f.v3
					}
				}
			}
		}
	}

	// clip Z far
	faceCount = fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		if f.visible {
			if f.z1 > farZ {
				if f.z2 > farZ {
					if f.z3 > farZ {
						f.visible = false
					} else {
						fakt := (farZ - f.z1) / (f.z1 - f.z3)
						f.x1 += (f.x1 - f.x3) * fakt
						f.y1 += (f.y1 - f.y3) * fakt
						f.z1 = farZ
						f.u1 += (f.u1 - f.u3) * fakt
						f.v1 += (f.v1 - f.v3) * fakt
						fakt = (farZ - f.z2) / (f.z2 - f.z3)
						f.x2 += (f.x2 - f.x3) * fakt
						f.y2 += (f.y2 - f.y3) * fakt
						f.z2 = farZ
						f.u2 += (f.u2 - f.u3) * fakt
						f.v2 += (f.v2 - f.v3) * fakt
					}
				} else {
					if f.z3 > farZ {
						fakt := (farZ - f.z1) / (f.z1 - f.z2)
						f.x1 += (f.x1 - f.x2) * fakt
						f.y1 += (f.y1 - f.y2) * fakt
						f.z1 = farZ
						f.u1 += (f.u1 - f.u2) * fakt
						f.v1 += (f.v1 - f.v2) * fakt
						fakt = (farZ - f.z3) / (f.z3 - f.z2)
						f.x3 += (f.x3 - f.x2) * fakt
						f.y3 += (f.y3 - f.y2) * fakt
						f.z3 = farZ
						f.u3 += (f.u3 - f.u2) * fakt
						f.v3 += (f.v3 - f.v2) * fakt
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z1 - farZ) / (f.z1 - f.z3)
						f.x1 += (f.x3 - f.x1) * fakt
						f.y1 += (f.y3 - f.y1) * fakt
						f.z1 = farZ
						f.u1 += (f.u3 - f.u1) * fakt
						f.v1 += (f.v3 - f.v1) * fakt
						fakt = (f1.z1 - farZ) / (f1.z1 - f1.z2)
						f1.x1 += (f1.x2 - f1.x1) * fakt
						f1.y1 += (f1.y2 - f1.y1) * fakt
						f1.z1 = farZ
						f1.u1 += (f1.u2 - f1.u1) * fakt
						f1.v1 += (f1.v2 - f1.v1) * fakt
						f1.x3 = f.x1
						f1.y3 = f.y1
						f1.z3 = farZ
						f1.u3 = f.u1
						f1.v3 = f.v1
					}
				}
			} else {
				if f.z2 > farZ {
					if f.z3 > farZ {
						fakt := (farZ - f.z2) / (f.z2 - f.z1)
						f.x2 += (f.x2 - f.x1) * fakt
						f.y2 += (f.y2 - f.y1) * fakt
						f.z2 = farZ
						f.u2 += (f.u2 - f.u1) * fakt
						f.v2 += (f.v2 - f.v1) * fakt
						fakt = (farZ - f.z3) / (f.z3 - f.z1)
						f.x3 += (f.x3 - f.x1) * fakt
						f.y3 += (f.y3 - f.y1) * fakt
						f.z3 = farZ
						f.u3 += (f.u3 - f.u1) * fakt
						f.v3 += (f.v3 - f.v1) * fakt
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - farZ) / (f.z2 - f.z3)
						f.x2 += (f.x3 - f.x2) * fakt
						f.y2 += (f.y3 - f.y2) * fakt
						f.z2 = farZ
						f.u2 += (f.u3 - f.u2) * fakt
						f.v2 += (f.v3 - f.v2) * fakt
						fakt = (f1.z2 - farZ) / (f1.z2 - f1.z1)
						f1.x2 += (f1.x1 - f1.x2) * fakt
						f1.y2 += (f1.y1 - f1.y2) * fakt
						f1.z2 = farZ
						f1.u2 += (f1.u1 - f1.u2) * fakt
						f1.v2 += (f1.v1 - f1.v2) * fakt
						f1.x3 = f.x2
						f1.y3 = f.y2
						f1.z3 = farZ
						f1.u3 = f.u2
						f1.v3 = f.v2
					}
				} else {
					if f.z3 > farZ {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z3 - farZ) / (f.z3 - f.z2)
						f.x3 += (f.x2 - f.x3) * fakt
						f.y3 += (f.y2 - f.y3) * fakt
						f.z3 = farZ
						f.u3 += (f.u2 - f.u3) * fakt
						f.v3 += (f.v2 - f.v3) * fakt
						fakt = (f1.z3 - farZ) / (f1.z3 - f1.z1)
						f1.x3 += (f1.x1 - f1.x3) * fakt
						f1.y3 += (f1.y1 - f1.y3) * fakt
						f1.z3 = farZ
						f1.u3 += (f1.u1 - f1.u3) * fakt
						f1.v3 += (f1.v1 - f1.v3) * fakt
						f1.x2 = f.x3
						f1.y2 = f.y3
						f1.z2 = farZ
						f1.u2 = f.u3
						f1.v2 = f.v3
					}
				}
			}
		}
	}

	// clip X left
	faceCount = fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		if f.visible {
			if f.x1 < -f.z1 {
				if f.x2 < -f.z2 {
					if f.x3 < -f.z3 {
						f.visible = false
					} else {
						fakt := (f.z3 - f.z1)
						skal := (-f.z1 - f.x1) / (f.x3 - f.x1 + fakt)
						neu := skal * fakt
						f.y1 += (f.y3 - f.y1) * skal
						f.u1 += (f.u3 - f.u1) * skal
						f.v1 += (f.v3 - f.v1) * skal
						f.z1 += neu
						f.x1 = -f.z1
						fakt = (f.z3 - f.z2)
						skal = (-f.z2 - f.x2) / (f.x3 - f.x2 + fakt)
						neu = skal * fakt
						f.y2 += (f.y3 - f.y2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.x2 = -f.z2
					}
				} else {
					if f.x3 < -f.z3 {
						fakt := (f.z2 - f.z1)
						skal := (-f.z1 - f.x1) / (f.x2 - f.x1 + fakt)
						neu := skal * fakt
						f.y1 += (f.y2 - f.y1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.x1 = -f.z1
						fakt = (f.z2 - f.z3)
						skal = (-f.z3 - f.x3) / (f.x2 - f.x3 + fakt)
						neu = skal * fakt
						f.y3 += (f.y2 - f.y3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.x3 = -f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z1)
						skal := (-f.z1 - f.x1) / (f.x2 - f.x1 + fakt)
						neu := skal * fakt
						f.y1 += (f.y2 - f.y1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.x1 = -f.z1
						fakt = (f1.z3 - f1.z1)
						skal = (-f1.z1 - f1.x1) / (f1.x3 - f1.x1 + fakt)
						neu = skal * fakt
						f1.y1 += (f1.y3 - f1.y1) * skal
						f1.u1 += (f1.u3 - f1.u1) * skal
						f1.v1 += (f1.v3 - f1.v1) * skal
						f1.z1 += neu
						f1.x1 = -f1.z1
						f1.x2 = f.x1
						f1.y2 = f.y1
						f1.z2 = f.z1
						f1.u2 = f.u1
						f1.v2 = f.v1
					}
				}
			} else {
				if f.x2 < -f.z2 {
					if f.x3 < -f.z3 {
						fakt := (f.z1 - f.z2)
						skal := (-f.z2 - f.x2) / (f.x1 - f.x2 + fakt)
						neu := skal * fakt
						f.y2 += (f.y1 - f.y2) * skal
						f.u2 += (f.u1 - f.u2) * skal
						f.v2 += (f.v1 - f.v2) * skal
						f.z2 += neu
						f.x2 = -f.z2
						fakt = (f.z1 - f.z3)
						skal = (-f.z3 - f.x3) / (f.x1 - f.x3 + fakt)
						neu = skal * fakt
						f.y3 += (f.y1 - f.y3) * skal
						f.u3 += (f.u1 - f.u3) * skal
						f.v3 += (f.v1 - f.v3) * skal
						f.z3 += neu
						f.x3 = -f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z3 - f.z2)
						skal := (-f.z2 - f.x2) / (f.x3 - f.x2 + fakt)
						neu := skal * fakt
						f.y2 += (f.y3 - f.y2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.x2 = -f.z2
						fakt = (f1.z1 - f1.z2)
						skal = (-f1.z2 - f1.x2) / (f1.x1 - f1.x2 + fakt)
						neu = skal * fakt
						f1.y2 += (f1.y1 - f1.y2) * skal
						f1.u2 += (f1.u1 - f1.u2) * skal
						f1.v2 += (f1.v1 - f1.v2) * skal
						f1.z2 += neu
						f1.x2 = -f1.z2
						f1.x3 = f.x2
						f1.y3 = f.y2
						f1.z3 = f.z2
						f1.u3 = f.u2
						f1.v3 = f.v2
					}
				} else {
					if f.x3 < -f.z3 {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z3)
						skal := (-f.z3 - f.x3) / (f.x2 - f.x3 + fakt)
						neu := skal * fakt
						f.y3 += (f.y2 - f.y3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.x3 = -f.z3
						fakt = (f1.z1 - f1.z3)
						skal = (-f1.z3 - f1.x3) / (f1.x1 - f1.x3 + fakt)
						neu = skal * fakt
						f1.y3 += (f1.y1 - f1.y3) * skal
						f1.u3 += (f1.u1 - f1.u3) * skal
						f1.v3 += (f1.v1 - f1.v3) * skal
						f1.z3 += neu
						f1.x3 = -f1.z3
						f1.x2 = f.x3
						f1.y2 = f.y3
						f1.z2 = f.z3
						f1.u2 = f.u3
						f1.v2 = f.v3
					}
				}
			}
		}
	}

	// clip X right
	faceCount = fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		if f.visible {
			if f.x1 > f.z1 {
				if f.x2 > f.z2 {
					if f.x3 > f.z3 {
						f.visible = false
					} else {
						fakt := (f.z3 - f.z1)
						skal := (f.z1 - f.x1) / (f.x3 - f.x1 - fakt)
						neu := skal * fakt
						f.y1 += (f.y3 - f.y1) * skal
						f.u1 += (f.u3 - f.u1) * skal
						f.v1 += (f.v3 - f.v1) * skal
						f.z1 += neu
						f.x1 = f.z1
						fakt = (f.z3 - f.z2)
						skal = (f.z2 - f.x2) / (f.x3 - f.x2 - fakt)
						neu = skal * fakt
						f.y2 += (f.y3 - f.y2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.x2 = f.z2
					}
				} else {
					if f.x3 > f.z3 {
						fakt := (f.z2 - f.z1)
						skal := (f.z1 - f.x1) / (f.x2 - f.x1 - fakt)
						neu := skal * fakt
						f.y1 += (f.y2 - f.y1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.x1 = f.z1
						fakt = (f.z2 - f.z3)
						skal = (f.z3 - f.x3) / (f.x2 - f.x3 - fakt)
						neu = skal * fakt
						f.y3 += (f.y2 - f.y3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.x3 = f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z1)
						skal := (f.z1 - f.x1) / (f.x2 - f.x1 - fakt)
						neu := skal * fakt
						f.y1 += (f.y2 - f.y1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.x1 = f.z1
						fakt = (f1.z3 - f1.z1)
						skal = (f1.z1 - f1.x1) / (f1.x3 - f1.x1 - fakt)
						neu = skal * fakt
						f1.y1 += (f1.y3 - f1.y1) * skal
						f1.u1 += (f1.u3 - f1.u1) * skal
						f1.v1 += (f1.v3 - f1.v1) * skal
						f1.z1 += neu
						f1.x1 = f1.z1
						f1.x2 = f.x1
						f1.y2 = f.y1
						f1.z2 = f.z1
						f1.u2 = f.u1
						f1.v2 = f.v1
					}
				}
			} else {
				if f.x2 > f.z2 {
					if f.x3 > f.z3 {
						fakt := (f.z1 - f.z2)
						skal := (f.z2 - f.x2) / (f.x1 - f.x2 - fakt)
						neu := skal * fakt
						f.y2 += (f.y1 - f.y2) * skal
						f.u2 += (f.u1 - f.u2) * skal
						f.v2 += (f.v1 - f.v2) * skal
						f.z2 += neu
						f.x2 = f.z2
						fakt = (f.z1 - f.z3)
						skal = (f.z3 - f.x3) / (f.x1 - f.x3 - fakt)
						neu = skal * fakt
						f.y3 += (f.y1 - f.y3) * skal
						f.u3 += (f.u1 - f.u3) * skal
						f.v3 += (f.v1 - f.v3) * skal
						f.z3 += neu
						f.x3 = f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z3 - f.z2)
						skal := (f.z2 - f.x2) / (f.x3 - f.x2 - fakt)
						neu := skal * fakt
						f.y2 += (f.y3 - f.y2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.x2 = f.z2
						fakt = (f1.z1 - f1.z2)
						skal = (f1.z2 - f1.x2) / (f1.x1 - f1.x2 - fakt)
						neu = skal * fakt
						f1.y2 += (f1.y1 - f1.y2) * skal
						f1.u2 += (f1.u1 - f1.u2) * skal
						f1.v2 += (f1.v1 - f1.v2) * skal
						f1.z2 += neu
						f1.x2 = f1.z2
						f1.x3 = f.x2
						f1.y3 = f.y2
						f1.z3 = f.z2
						f1.u3 = f.u2
						f1.v3 = f.v2
					}
				} else {
					if f.x3 > f.z3 {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z3)
						skal := (f.z3 - f.x3) / (f.x2 - f.x3 - fakt)
						neu := fakt * skal
						f.y3 += (f.y2 - f.y3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.x3 = f.z3
						fakt = (f1.z1 - f1.z3)
						skal = (f1.z3 - f1.x3) / (f1.x1 - f1.x3 - fakt)
						neu = fakt * skal
						f1.y3 += (f1.y1 - f1.y3) * skal
						f1.u3 += (f1.u1 - f1.u3) * skal
						f1.v3 += (f1.v1 - f1.v3) * skal
						f1.z3 += neu
						f1.x3 = f1.z3
						f1.x2 = f.x3
						f1.y2 = f.y3
						f1.z2 = f.z3
						f1.u2 = f.u3
						f1.v2 = f.v3
					}
				}
			}
		}
	}

	// clip Y bottom
	faceCount = fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		if f.visible {
			if f.y1 < -f.z1 {
				if f.y2 < -f.z2 {
					if f.y3 < -f.z3 {
						f.visible = false
					} else {
						fakt := (f.z3 - f.z1)
						skal := (-f.z1 - f.y1) / (f.y3 - f.y1 + fakt)
						neu := fakt * skal
						f.x1 += (f.x3 - f.x1) * skal
						f.u1 += (f.u3 - f.u1) * skal
						f.v1 += (f.v3 - f.v1) * skal
						f.z1 += neu
						f.y1 = -f.z1
						fakt = (f.z3 - f.z2)
						skal = (-f.z2 - f.y2) / (f.y3 - f.y2 + fakt)
						neu = fakt * skal
						f.x2 += (f.x3 - f.x2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.y2 = -f.z2
					}
				} else {
					if f.y3 < -f.z3 {
						fakt := (f.z2 - f.z1)
						skal := (-f.z1 - f.y1) / (f.y2 - f.y1 + fakt)
						neu := fakt * skal
						f.x1 += (f.x2 - f.x1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.y1 = -f.z1
						fakt = (f.z2 - f.z3)
						skal = (-f.z3 - f.y3) / (f.y2 - f.y3 + fakt)
						neu = fakt * skal
						f.x3 += (f.x2 - f.x3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.y3 = -f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z1)
						skal := (-f.z1 - f.y1) / (f.y2 - f.y1 + fakt)
						neu := fakt * skal
						f.x1 += (f.x2 - f.x1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.y1 = -f.z1
						fakt = (f1.z3 - f1.z1)
						skal = (-f1.z1 - f1.y1) / (f1.y3 - f1.y1 + fakt)
						neu = fakt * skal
						f1.x1 += (f1.x3 - f1.x1) * skal
						f1.u1 += (f1.u3 - f1.u1) * skal
						f1.v1 += (f1.v3 - f1.v1) * skal
						f1.z1 += neu
						f1.y1 = -f1.z1
						f1.y2 = f.y1
						f1.x2 = f.x1
						f1.z2 = f.z1
						f1.v2 = f.v1
						f1.u2 = f.u1
					}
				}
			} else {
				if f.y2 < -f.z2 {
					if f.y3 < -f.z3 {
						fakt := (f.z1 - f.z2)
						skal := (-f.z2 - f.y2) / (f.y1 - f.y2 + fakt)
						neu := fakt * skal
						f.x2 += (f.x1 - f.x2) * skal
						f.u2 += (f.u1 - f.u2) * skal
						f.v2 += (f.v1 - f.v2) * skal
						f.z2 += neu
						f.y2 = -f.z2
						fakt = (f.z1 - f.z3)
						skal = (-f.z3 - f.y3) / (f.y1 - f.y3 + fakt)
						neu = fakt * skal
						f.x3 += (f.x1 - f.x3) * skal
						f.u3 += (f.u1 - f.u3) * skal
						f.v3 += (f.v1 - f.v3) * skal
						f.z3 += neu
						f.y3 = -f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z3 - f.z2)
						skal := (-f.z2 - f.y2) / (f.y3 - f.y2 + fakt)
						neu := fakt * skal
						f.x2 += (f.x3 - f.x2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.y2 = -f.z2
						fakt = (f1.z1 - f1.z2)
						skal = (-f1.z2 - f1.y2) / (f1.y1 - f1.y2 + fakt)
						neu = fakt * skal
						f1.x2 += (f1.x1 - f1.x2) * skal
						f1.u2 += (f1.u1 - f1.u2) * skal
						f1.v2 += (f1.v1 - f1.v2) * skal
						f1.z2 += neu
						f1.y2 = -f1.z2
						f1.y3 = f.y2
						f1.x3 = f.x2
						f1.z3 = f.z2
						f1.v3 = f.v2
						f1.u3 = f.u2
					}
				} else {
					if f.y3 < -f.z3 {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z3)
						skal := (-f.z3 - f.y3) / (f.y2 - f.y3 + fakt)
						neu := fakt * skal
						f.x3 += (f.x2 - f.x3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.y3 = -f.z3
						fakt = (f1.z1 - f1.z3)
						skal = (-f1.z3 - f1.y3) / (f1.y1 - f1.y3 + fakt)
						neu = fakt * skal
						f1.x3 += (f1.x1 - f1.x3) * skal
						f1.u3 += (f1.u1 - f1.u3) * skal
						f1.v3 += (f1.v1 - f1.v3) * skal
						f1.z3 += neu
						f1.y3 = -f1.z3
						f1.y2 = f.y3
						f1.x2 = f.x3
						f1.z2 = f.z3
						f1.v2 = f.v3
						f1.u2 = f.u3
					}
				}
			}
		}
	}

	// clip Y top
	faceCount = fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		if f.visible {
			if f.y1 > f.z1 {
				if f.y2 > f.z2 {
					if f.y3 > f.z3 {
						f.visible = false
					} else {
						fakt := (f.z3 - f.z1)
						skal := (f.z1 - f.y1) / (f.y3 - f.y1 - fakt)
						neu := fakt * skal
						f.x1 += (f.x3 - f.x1) * skal
						f.u1 += (f.u3 - f.u1) * skal
						f.v1 += (f.v3 - f.v1) * skal
						f.z1 += neu
						f.y1 = f.z1
						fakt = (f.z3 - f.z2)
						skal = (f.z2 - f.y2) / (f.y3 - f.y2 - fakt)
						neu = fakt * skal
						f.x2 += (f.x3 - f.x2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.y2 = f.z2
					}
				} else {
					if f.y3 > f.z3 {
						fakt := (f.z2 - f.z1)
						skal := (f.z1 - f.y1) / (f.y2 - f.y1 - fakt)
						neu := fakt * skal
						f.x1 += (f.x2 - f.x1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.y1 = f.z1
						fakt = (f.z2 - f.z3)
						skal = (f.z3 - f.y3) / (f.y2 - f.y3 - fakt)
						neu = fakt * skal
						f.x3 += (f.x2 - f.x3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.y3 = f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z1)
						skal := (f.z1 - f.y1) / (f.y2 - f.y1 - fakt)
						neu := fakt * skal
						f.x1 += (f.x2 - f.x1) * skal
						f.u1 += (f.u2 - f.u1) * skal
						f.v1 += (f.v2 - f.v1) * skal
						f.z1 += neu
						f.y1 = f.z1
						fakt = (f1.z3 - f1.z1)
						skal = (f1.z1 - f1.y1) / (f1.y3 - f1.y1 - fakt)
						neu = fakt * skal
						f1.x1 += (f1.x3 - f1.x1) * skal
						f1.u1 += (f1.u3 - f1.u1) * skal
						f1.v1 += (f1.v3 - f1.v1) * skal
						f1.z1 += neu
						f1.y1 = f1.z1
						f1.y2 = f.y1
						f1.x2 = f.x1
						f1.z2 = f.z1
						f1.v2 = f.v1
						f1.u2 = f.u1
					}
				}
			} else {
				if f.y2 > f.z2 {
					if f.y3 > f.z3 {
						fakt := (f.z1 - f.z2)
						skal := (f.z2 - f.y2) / (f.y1 - f.y2 - fakt)
						neu := fakt * skal
						f.x2 += (f.x1 - f.x2) * skal
						f.u2 += (f.u1 - f.u2) * skal
						f.v2 += (f.v1 - f.v2) * skal
						f.z2 += neu
						f.y2 = f.z2
						fakt = (f.z1 - f.z3)
						skal = (f.z3 - f.y3) / (f.y1 - f.y3 - fakt)
						neu = fakt * skal
						f.x3 += (f.x1 - f.x3) * skal
						f.u3 += (f.u1 - f.u3) * skal
						f.v3 += (f.v1 - f.v3) * skal
						f.z3 += neu
						f.y3 = f.z3
					} else {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z3 - f.z2)
						skal := (f.z2 - f.y2) / (f.y3 - f.y2 - fakt)
						neu := fakt * skal
						f.x2 += (f.x3 - f.x2) * skal
						f.u2 += (f.u3 - f.u2) * skal
						f.v2 += (f.v3 - f.v2) * skal
						f.z2 += neu
						f.y2 = f.z2
						fakt = (f1.z1 - f1.z2)
						skal = (f1.z2 - f1.y2) / (f1.y1 - f1.y2 - fakt)
						neu = fakt * skal
						f1.x2 += (f1.x1 - f1.x2) * skal
						f1.u2 += (f1.u1 - f1.u2) * skal
						f1.v2 += (f1.v1 - f1.v2) * skal
						f1.z2 += neu
						f1.y2 = f1.z2
						f1.y3 = f.y2
						f1.x3 = f.x2
						f1.z3 = f.z2
						f1.v3 = f.v2
						f1.u3 = f.u2
					}
				} else {
					if f.y3 > f.z3 {
						f1 := fl.nextFace()
						if f1 == nil {
							continue
						}
						*f1 = *f
						fakt := (f.z2 - f.z3)
						skal := (f.z3 - f.y3) / (f.y2 - f.y3 - fakt)
						neu := fakt * skal
						f.x3 += (f.x2 - f.x3) * skal
						f.u3 += (f.u2 - f.u3) * skal
						f.v3 += (f.v2 - f.v3) * skal
						f.z3 += neu
						f.y3 = f.z3
						fakt = (f1.z1 - f1.z3)
						skal = (f1.z3 - f1.y3) / (f1.y1 - f1.y3 - fakt)
						neu = fakt * skal
						f1.x3 += (f1.x1 - f1.x3) * skal
						f1.u3 += (f1.u1 - f1.u3) * skal
						f1.v3 += (f1.v1 - f1.v3) * skal
						f1.z3 += neu
						f1.y3 = f1.z3
						f1.y2 = f.y3
						f1.x2 = f.x3
						f1.z2 = f.z3
						f1.v2 = f.v3
						f1.u2 = f.u3
					}
				}
			}
		}
	}

	// remove not visible and project
	d_anf := 0
	faceCount = fl.count
	for t := 0; t < faceCount; t++ {
		f := dl[t]
		if f.visible {
			if f.z1 < f.z2 {
				f.minZ = f.z1
				f.maxZ = f.z2
			} else {
				f.minZ = f.z2
				f.maxZ = f.z1
			}
			if f.z3 < f.minZ {
				f.minZ = f.z3
			}
			if f.z3 > f.maxZ {
				f.maxZ = f.z3
			}

			z1 := 1.0 / f.z1
			z2 := 1.0 / f.z2
			z3 := 1.0 / f.z3
			x1 := f.x1 * z1
			y1 := f.y1 * z1
			x2 := f.x2 * z2
			y2 := f.y2 * z2
			x3 := f.x3 * z3
			y3 := f.y3 * z3

			f.x1 = (320.0 + 319.0*x1)
			f.y1 = (240.0 + 239.0*y1)
			f.x2 = (320.0 + 319.0*x2)
			f.y2 = (240.0 + 239.0*y2)
			f.x3 = (320.0 + 319.0*x3)
			f.y3 = (240.0 + 239.0*y3)

			if f.x1 < f.x2 {
				f.minx = f.x1
				f.maxx = f.x2
			} else {
				f.minx = f.x2
				f.maxx = f.x1
			}
			if f.x3 < f.minx {
				f.minx = f.x3
			}
			if f.x3 > f.maxx {
				f.maxx = f.x3
			}

			if f.y1 < f.y2 {
				f.miny = f.y1
				f.maxy = f.y2
			} else {
				f.miny = f.y2
				f.maxy = f.y1
			}
			if f.y3 < f.miny {
				f.miny = f.y3
			}
			if f.y3 > f.maxy {
				f.maxy = f.y3
			}

			dl[d_anf] = f
			d_anf++
		}
	}
	fl.count = d_anf

	if fl.count > 0 {
		if mz {
			fl.sortMin()
		} else {
			fl.sortMax()
		}
	}
}

func (fl *FaceList) sortMax() {
	sort.Slice(fl.dl[:fl.count], func(i, j int) bool {
		return fl.dl[i].maxZ > fl.dl[j].maxZ
	})
}

func (fl *FaceList) sortMin() {
	sort.Slice(fl.dl[:fl.count], func(i, j int) bool {
		return fl.dl[i].minZ > fl.dl[j].minZ
	})
}

func (fl *FaceList) nextFace() *DrawFaceObj {
	if fl.scratchCount >= len(fl.scratchFaces) || fl.count >= len(fl.dl) {
		return nil
	}
	f1 := &fl.scratchFaces[fl.scratchCount]
	fl.scratchCount++
	fl.dl[fl.count] = f1
	fl.count++
	return f1
}
