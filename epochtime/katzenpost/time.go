// timer.go - Reunion Katzenpost epoch timer.
// Copyright (C) 2020  David Stainton.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package time provides the Reunion protocol Katzenpost epoch timer.
package time

import (
	"time"

	"github.com/katzenpost/epochtime"
)

// EpochTimer provides the Epochtimer interface for the
// Katzenpost epoch timer.
type EpochTimer struct{}

func (t *EpochTimer) Now() (current uint64, elapsed, till time.Duration) {
	return epochtime.Now()
}
