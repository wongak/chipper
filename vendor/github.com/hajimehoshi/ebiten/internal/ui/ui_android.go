// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

/*

#include <jni.h>
#include <stdlib.h>

// Basically same as `getResources().getDisplayMetrics().density`;
static float deviceScale(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx) {
  JavaVM* vm = (JavaVM*)java_vm;
  JNIEnv* env = (JNIEnv*)jni_env;
  jobject context = (jobject)ctx;

  const jclass android_content_ContextWrapper =
      (*env)->FindClass(env, "android/content/ContextWrapper");
  const jclass android_content_res_Resources =
      (*env)->FindClass(env, "android/content/res/Resources");
  const jclass android_util_DisplayMetrics =
      (*env)->FindClass(env, "android/util/DisplayMetrics");

  const jobject resources =
      (*env)->CallObjectMethod(
          env, context,
          (*env)->GetMethodID(env, android_content_ContextWrapper, "getResources", "()Landroid/content/res/Resources;"));
  const jobject displayMetrics =
      (*env)->CallObjectMethod(
          env, resources,
          (*env)->GetMethodID(env, android_content_res_Resources, "getDisplayMetrics", "()Landroid/util/DisplayMetrics;"));
  const float density =
      (*env)->GetFloatField(
          env, displayMetrics,
          (*env)->GetFieldID(env, android_util_DisplayMetrics, "density", "F"));
  return density;
}

*/
import "C"

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/internal/jni"
)

var (
	androidDeviceScale = 0.0
)

func deviceScale() float64 {
	if 0 < androidDeviceScale {
		return androidDeviceScale
	}
	if err := jni.RunOnJVM(func(vm, env, ctx uintptr) error {
		androidDeviceScale = float64(C.deviceScale(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx)))
		return nil
	}); err != nil {
		panic(fmt.Sprintf("ui: error %v", err))
	}
	return androidDeviceScale
}
