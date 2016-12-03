package commune

func (p Freshness0) Len() int {             return len(p)}
func (p Freshness0) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness0) Less(i, j int) bool {   return value(0, posts[p[i]]) > value(0, posts[p[j]])}

func (p Freshness1) Len() int {             return len(p)}
func (p Freshness1) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness1) Less(i, j int) bool {   return value(0.1, posts[p[i]]) > value(0.1, posts[p[j]])}

func (p Freshness2) Len() int {             return len(p)}
func (p Freshness2) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness2) Less(i, j int) bool {   return value(0.2, posts[p[i]]) > value(0.2, posts[p[j]])}

func (p Freshness3) Len() int {             return len(p)}
func (p Freshness3) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness3) Less(i, j int) bool {   return value(0.5, posts[p[i]]) > value(0.5, posts[p[j]])}

func (p Freshness4) Len() int {             return len(p)}
func (p Freshness4) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness4) Less(i, j int) bool {   return value(1, posts[p[i]]) > value(1, posts[p[j]])}

func (p Freshness5) Len() int {             return len(p)}
func (p Freshness5) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness5) Less(i, j int) bool {   return value(2, posts[p[i]]) > value(2, posts[p[j]])}

func (p Freshness6) Len() int {             return len(p)}
func (p Freshness6) Swap(i, j int) {        p[i], p[j] = p[j], p[i]}
func (p Freshness6) Less(i, j int) bool {   return value(5, posts[p[i]]) > value(5, posts[p[j]])}
