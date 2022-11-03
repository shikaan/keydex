package utils

func Merge[K string, V interface{}] (m1, m2 map[K]V) map[K]V {
  m3 := map[K]V{}

  for k,v := range m1 {
    m3[k] = v
  }

  for k,v := range m2 {
    m3[k] = v
  }

  return m3
}

func Max(a, b int) int {
  if a < b {
    return b
  }
  return a
}

func Min(a,b int) int {
  if a > b {
    return b
  }
  return a
}
