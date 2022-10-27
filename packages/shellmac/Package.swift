// swift-tools-version: 5.4
import PackageDescription

let package = Package(
    name: "KoiShell",
    platforms: [
        .macOS(.v11)
    ],
    dependencies: [
        .package(url: "https://github.com/kylehickinson/SwiftUI-WebView", .upToNextMinor(from: "0.3.0")),
    ],
    targets: [
        .executableTarget(
            name: "KoiShell",
            dependencies: [
                .product(name: "WebView", package: "SwiftUI-WebView"),
            ]),
    ]
)
