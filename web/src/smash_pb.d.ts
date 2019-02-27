// package: proto
// file: smash.proto

import * as jspb from 'google-protobuf';

export class RunRequest extends jspb.Message {
  getCell(): number;
  setCell(value: number): void;

  getCwd(): string;
  setCwd(value: string): void;

  clearArgvList(): void;
  getArgvList(): Array<string>;
  setArgvList(value: Array<string>): void;
  addArgv(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: RunRequest
  ): RunRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: RunRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): RunRequest;
  static deserializeBinaryFromReader(
    message: RunRequest,
    reader: jspb.BinaryReader
  ): RunRequest;
}

export namespace RunRequest {
  export type AsObject = {
    cell: number;
    cwd: string;
    argvList: Array<string>;
  };
}

export class TermAttr extends jspb.Message {
  getFg(): number;
  setFg(value: number): void;

  getBg(): number;
  setBg(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TermAttr.AsObject;
  static toObject(includeInstance: boolean, msg: TermAttr): TermAttr.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: TermAttr,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): TermAttr;
  static deserializeBinaryFromReader(
    message: TermAttr,
    reader: jspb.BinaryReader
  ): TermAttr;
}

export namespace TermAttr {
  export type AsObject = {
    fg: number;
    bg: number;
  };
}

export class TermText extends jspb.Message {
  getRow(): number;
  setRow(value: number): void;

  clearSpansList(): void;
  getSpansList(): Array<TermText.Span>;
  setSpansList(value: Array<TermText.Span>): void;
  addSpans(value?: TermText.Span, index?: number): TermText.Span;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TermText.AsObject;
  static toObject(includeInstance: boolean, msg: TermText): TermText.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: TermText,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): TermText;
  static deserializeBinaryFromReader(
    message: TermText,
    reader: jspb.BinaryReader
  ): TermText;
}

export namespace TermText {
  export type AsObject = {
    row: number;
    spansList: Array<TermText.Span.AsObject>;
  };

  export class Span extends jspb.Message {
    hasAttr(): boolean;
    clearAttr(): void;
    getAttr(): TermAttr | undefined;
    setAttr(value?: TermAttr): void;

    getText(): string;
    setText(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Span.AsObject;
    static toObject(includeInstance: boolean, msg: Span): Span.AsObject;
    static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
    static extensionsBinary: {
      [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
    };
    static serializeBinaryToWriter(
      message: Span,
      writer: jspb.BinaryWriter
    ): void;
    static deserializeBinary(bytes: Uint8Array): Span;
    static deserializeBinaryFromReader(
      message: Span,
      reader: jspb.BinaryReader
    ): Span;
  }

  export namespace Span {
    export type AsObject = {
      attr?: TermAttr.AsObject;
      text: string;
    };
  }
}

export class TermState extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TermState.AsObject;
  static toObject(includeInstance: boolean, msg: TermState): TermState.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: TermState,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): TermState;
  static deserializeBinaryFromReader(
    message: TermState,
    reader: jspb.BinaryReader
  ): TermState;
}

export namespace TermState {
  export type AsObject = {};
}

export class Output extends jspb.Message {
  getCell(): number;
  setCell(value: number): void;

  hasError(): boolean;
  clearError(): void;
  getError(): string;
  setError(value: string): void;

  hasText(): boolean;
  clearText(): void;
  getText(): TermText | undefined;
  setText(value?: TermText): void;

  hasExitCode(): boolean;
  clearExitCode(): void;
  getExitCode(): number;
  setExitCode(value: number): void;

  getOutputCase(): Output.OutputCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Output.AsObject;
  static toObject(includeInstance: boolean, msg: Output): Output.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Output,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Output;
  static deserializeBinaryFromReader(
    message: Output,
    reader: jspb.BinaryReader
  ): Output;
}

export namespace Output {
  export type AsObject = {
    cell: number;
    error: string;
    text?: TermText.AsObject;
    exitCode: number;
  };

  export enum OutputCase {
    OUTPUT_NOT_SET = 0,
    ERROR = 2,
    TEXT = 3,
    EXIT_CODE = 4
  }
}

export class ServerMsg extends jspb.Message {
  hasOutput(): boolean;
  clearOutput(): void;
  getOutput(): Output | undefined;
  setOutput(value?: Output): void;

  getMsgCase(): ServerMsg.MsgCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerMsg.AsObject;
  static toObject(includeInstance: boolean, msg: ServerMsg): ServerMsg.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ServerMsg,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ServerMsg;
  static deserializeBinaryFromReader(
    message: ServerMsg,
    reader: jspb.BinaryReader
  ): ServerMsg;
}

export namespace ServerMsg {
  export type AsObject = {
    output?: Output.AsObject;
  };

  export enum MsgCase {
    MSG_NOT_SET = 0,
    OUTPUT = 1
  }
}