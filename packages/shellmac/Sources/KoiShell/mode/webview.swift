import AppKit
import SwiftUI
import WebView
import WebKit

struct KSWebViewOutput: Codable {}

func ksWebView(_ arg: [String: Any]) {
    guard let url = arg["url"] as? String else {
        log("Failed to parse url.")
        exit(EXIT_FAILURE)
    }

    guard let name = arg["name"] as? String else {
        log("Failed to parse name.")
        exit(EXIT_FAILURE)
    }

    let delegate = KSWebViewDelegate(url, name)
    NSApp.delegate = delegate
    NSApp.run()
}

class KSWebViewDelegate: NSObject, NSApplicationDelegate, NSWindowDelegate, WKScriptMessageHandler {
    var window: NSWindow!
    var hostingView: NSView?
    let name: String
    let url: String
    var appearance: NSAppearance!
    var initAppearance: NSAppearance!

    init(_ url: String, _ name: String) {
        self.url = url
        self.name = name
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        let contentView = KSWebView(url, self)

        window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 1366, height: 768),
            styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView],
            backing: .buffered,
            defer: false
        )
        appearance = window.appearance
        initAppearance = window.appearance

        window.title = "\(self.name) - Koishi"
        window.titlebarAppearsTransparent = true
        // window.titleVisibility = .hidden
        // window.appearance = NSAppearance(named: NSAppearance.Name.vibrantLight)
        // window.isMovableByWindowBackground = true
        window.center()
        window.setFrameAutosaveName("KSWebView")
        hostingView = NSHostingView(rootView: contentView)
        window.contentView = hostingView
        window.delegate = self

        window.makeKeyAndOrderFront(nil)
        NSApp.activate(ignoringOtherApps: true)
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool {
        NSApplication.shared.terminate(self)
        return true
    }

    func windowWillClose(_ notification: Notification) {
        print(try! JSONEncoder().encode(KSWebViewOutput()).base64EncodedString())
    }

    func userContentController(_ userContentController: WKUserContentController, didReceive message: WKScriptMessage) {
        guard let msg = message.body as? String else {
            log("Failed to parse shellmacHandler message \(message.body)")
            return
        }

        switch msg {
        case "TL":
            appearance = NSAppearance(named: NSAppearance.Name.vibrantLight)
        case "TD":
            appearance = NSAppearance(named: NSAppearance.Name.vibrantDark)
        case "TR":
            appearance = initAppearance
        default:
            break
        }

        setAppearance()
    }

    func setAppearance() {
        window.appearance = appearance
        window.appearanceSource.appearance = appearance
        window.invalidateShadow()
    }

    func windowDidBecomeKey(_ notification: Notification) {
        setAppearance()
    }

    func windowDidResignKey(_ notification: Notification) {
        setAppearance()
    }
}

struct KSWebView: View {
    @StateObject var webViewStore: WebViewStore
    private var url: String

    init(_ url: String, _ webViewDelegate: KSWebViewDelegate) {
        self.url = url

        let enhanceURL = Bundle.module.url(forResource: "userscript", withExtension: "js")!
        let enhanceData = try! Data(contentsOf: enhanceURL)
        let enhanceRaw = String(decoding: enhanceData, as: UTF8.self)
        let userScript = WKUserScript(
            source: enhanceRaw,
            injectionTime: .atDocumentEnd,
            forMainFrameOnly: true
        )
        let userContentController = WKUserContentController()
        userContentController.addUserScript(userScript)
        userContentController.add(webViewDelegate, name: "shellmacHandler")
        let configuration = WKWebViewConfiguration()
        configuration.userContentController = userContentController
        _webViewStore = StateObject(
            wrappedValue: WebViewStore(
                webView: WKWebView(
                    frame: .zero,
                    configuration: configuration
                )
            )
        )
    }

    var body: some View {
        ZStack {
            ZStack {}
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .visualEffectBackground()
                .ignoresSafeArea()

            ZStack {
                WebView(webView: webViewStore.webView)
                    .onAppear {
                        self.webViewStore.webView.setValue(false, forKey: "drawsBackground")
                        self.webViewStore.webView.configuration.preferences.setValue(true, forKey: "developerExtrasEnabled")
                        self.webViewStore.webView.load(URLRequest(url: URL(string: self.url)!))
                    }
                    .visualEffectBackground()
            }
        }
    }
}
