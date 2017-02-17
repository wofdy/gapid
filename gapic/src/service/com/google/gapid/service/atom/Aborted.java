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
package com.google.gapid.service.atom;

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

public final class Aborted implements BinaryObject {
  //<<<Start:Java.ClassBody:1>>>
  private boolean myIsAssert;
  private String myReason;

  // Constructs a default-initialized {@link Aborted}.
  public Aborted() {}


  public boolean getIsAssert() {
    return myIsAssert;
  }

  public Aborted setIsAssert(boolean v) {
    myIsAssert = v;
    return this;
  }

  public String getReason() {
    return myReason;
  }

  public Aborted setReason(String v) {
    myReason = v;
    return this;
  }

  @Override
  public BinaryClass klass() { return Klass.INSTANCE; }


  private static final Entity ENTITY = new Entity("atom", "Aborted", "", "");

  static {
    ENTITY.setFields(new Field[]{
      new Field("IsAssert", new Primitive("bool", Method.Bool)),
      new Field("Reason", new Primitive("string", Method.String)),
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
    public BinaryObject create() { return new Aborted(); }

    @Override
    public void encode(Encoder e, BinaryObject obj) throws IOException {
      Aborted o = (Aborted)obj;
      e.bool(o.myIsAssert);
      e.string(o.myReason);
    }

    @Override
    public void decode(Decoder d, BinaryObject obj) throws IOException {
      Aborted o = (Aborted)obj;
      o.myIsAssert = d.bool();
      o.myReason = d.string();
    }
    //<<<End:Java.KlassBody:2>>>
  }
}
