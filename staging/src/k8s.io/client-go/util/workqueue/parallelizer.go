/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package workqueue

import (
	"context"
	"sync"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

type DoWorkPieceFunc func(piece int)

// Parallelize is a very simple framework that allows for parallelizing
// N independent pieces of work.
//
// Deprecated: Use ParallelizeUntil instead.
func Parallelize(workers, pieces int, doWorkPiece DoWorkPieceFunc) {
	ParallelizeUntil(nil, workers, pieces, doWorkPiece)
}

// ParallelizeUntil is a framework that allows for parallelizing N
// independent pieces of work until done or the context is canceled.
func ParallelizeUntil(ctx context.Context, workers, pieces int, doWorkPiece DoWorkPieceFunc) {
	var stop <-chan struct{}
	if ctx != nil {
		stop = ctx.Done()
	}

    // 将pieces放入toProcess中
	toProcess := make(chan int, pieces)
	for i := 0; i < pieces; i++ {
		toProcess <- i

	}
	close(toProcess)

	if pieces < workers {
		workers = pieces
	}

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer utilruntime.HandleCrash()
			defer wg.Done()

			// 从toProcess读取piece索引, 并做出处理(toProcess是channel, 可以并发读)
			for piece := range toProcess {
				select {
				case <-stop:
					// 如果ctx结束, 那么这里返回, 不再处理piece
                    // 可以看出: 如果ctx通知ParallelizeUntil结束, 需要在一个doWorkPiece()完全处理完毕后
					return
				default:
					doWorkPiece(piece)
				}
			}
		}()
	}
	wg.Wait()
}
