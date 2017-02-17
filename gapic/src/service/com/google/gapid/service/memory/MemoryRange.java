/*
 * Copyright (C) 2017 Google Inc.
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
package com.google.gapid.service.memory;

import com.google.gapid.rpclib.binary.BinaryClass;
import com.google.gapid.rpclib.binary.BinaryObject;
import com.google.gapid.rpclib.binary.Decoder;
import com.google.gapid.rpclib.binary.Encoder;
import com.google.gapid.rpclib.binary.Namespace;
import com.google.gapid.rpclib.schema.Entity;
import com.google.gapid.rpclib.schema.Field;
import com.google.gapid.rpclib.schema.Method;
import com.google.gapid.rpclib.schema.Primitive;

import java.io.IOException;

public final class MemoryRange implements BinaryObject {
  public boolean contains(long value) {
    return value >= myBase && value < myBase + mySize;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;
    if (o == null || getClass() != o.getClass()) return false;
    MemoryRange that = (MemoryRange)o;
    if (myBase != that.myBase) return false;
    if (mySize != that.mySize) return false;
    return true;
  }

  @Override
  public int hashCode() {
    int result = (int)(myBase ^ (myBase >>> 32));
    result = 31 * result + (int)(mySize ^ (mySize >>> 32));
    return result;
  }

  @Override
  public String toString() {
    return "MemoryRange{base=" + myBase + ", size=" + mySize + '}';
  }

  //<<<Start:Java.ClassBody:1>>>
  private long myBase;
  private long mySize;

  // Constructs a default-initialized {@link MemoryRange}.
  public MemoryRange() {}


  public long getBase() {
    return myBase;
  }

  public MemoryRange setBase(long v) {
    myBase = v;
    return this;
  }

  public long getSize() {
    return mySize;
  }

  public MemoryRange setSize(long v) {
    mySize = v;
    return this;
  }

  @Override
  public BinaryClass klass() { return Klass.INSTANCE; }


  private static final Entity ENTITY = new Entity("memory", "Range", "", "");

  static {
    ENTITY.setFields(new Field[]{
      new Field("Base", new Primitive("uint64", Method.Uint64)),
      new Field("Size", new Primitive("uint64", Method.Uint64)),
    });
    Namespace.register(Klass.INSTANCE);
  }
  public static void register() {}
  //<<<End:Java.ClassBody:1>>>
  public enum Klass implements BinaryClass {
    //<<<Start:Java.KlassBody:2>>>
    INSTANCE;

    @Override
    public Entity entity() { return ENTITY; }

    @Override
    public BinaryObject create() { return new MemoryRange(); }

    @Override
    public void encode(Encoder e, BinaryObject obj) throws IOException {
      MemoryRange o = (MemoryRange)obj;
      e.uint64(o.myBase);
      e.uint64(o.mySize);
    }

    @Override
    public void decode(Decoder d, BinaryObject obj) throws IOException {
      MemoryRange o = (MemoryRange)obj;
      o.myBase = d.uint64();
      o.mySize = d.uint64();
    }
    //<<<End:Java.KlassBody:2>>>
  }
}
