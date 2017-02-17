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

var path = require('path');
var vscode = require('vscode');
var langClient = require('vscode-languageclient');

var LanguageClient = langClient.LanguageClient;

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
function activate(context) {
	let serverModule = path.join(context.extensionPath, 'bin', 'langsvr');
	let debugOptions = { execArgv: ["--debug"] };

	// If the extension is launched in debug mode then the debug server options are used
	// Otherwise the run options are used
	let serverOptions = {
		run: { command: serverModule, transport: langClient.stdio },
		debug: { command: serverModule, transport: langClient.stdio, options: debugOptions }
	}

	// Options to control the language client
	let clientOptions = {
		documentSelector: ['gfxapi'],
		synchronize: {
			// Synchronize the setting section 'gfxapi' to the server
			configurationSection: 'gfxapi',
			// Notify the server about file changes to .api files contained in the workspace
			fileEvents: vscode.workspace.createFileSystemWatcher('**/*.api')
		}
	}

	// Create the language client and start the client.
	let disposable = new LanguageClient('gfxapi', serverOptions, clientOptions).start();

	// Push the disposable to the context's subscriptions so that the
	// client can be deactivated on extension deactivation
	context.subscriptions.push(disposable);
}
exports.activate = activate;

// this method is called when your extension is deactivated
function deactivate() {
}
exports.deactivate = deactivate;