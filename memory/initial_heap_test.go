/*
 * Copyright 2026 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package memory_test

import (
	"testing"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/v4/memory"
	. "github.com/onsi/gomega"
)

func TestInitialHeap(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(memory.InitialHeap(memory.Kibi).String()).To(Equal("-Xms1K"))
	g.Expect(memory.IsInitialHeap(" -Xms1K ")).To(BeTrue())
	g.Expect(memory.IsInitialHeap("-Xmx1K")).To(BeFalse())
	g.Expect(memory.ParseInitialHeap(" -Xms1K ")).To(Equal(memory.InitialHeap(memory.Kibi)))
	_, err := memory.ParseInitialHeap("-Xmsinvalid")
	g.Expect(err).To(HaveOccurred())
	_, err = memory.ParseInitialHeap("-Xms999999999999999999999999")
	g.Expect(err).To(HaveOccurred())
}
