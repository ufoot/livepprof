// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it wil/l be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

type LP struct {
}

var _ Profiler = &LP{}

func New() *LP {
	// [TODO:ufoot] implement this
	return &LP{}
}

func (lp *LP) Close() {
	// [TODO:ufoot] implement this
}
