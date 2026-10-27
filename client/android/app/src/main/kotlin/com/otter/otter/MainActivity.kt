package com.otter.otter

import android.content.Intent
import android.net.Uri
import android.os.Build
import android.provider.Settings
import androidx.core.content.FileProvider
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import java.io.File

// In-app updates (DD-304): Dart downloads the APK to cacheDir/updates/update.apk
// and hands it to the system installer, which only accepts it when it is
// signed with this app's key. The file name is fixed here, never taken from Dart.
class MainActivity : FlutterActivity() {
    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, "mena/update").setMethodCallHandler { call, result ->
            when (call.method) {
                "dir" -> result.success(updatesDir().absolutePath)
                "install" -> result.success(install())
                else -> result.notImplemented()
            }
        }
    }

    private fun updatesDir() = File(cacheDir, "updates").apply { mkdirs() }

    // "permission" when installs from this app must be allowed first; the
    // settings page is opened and the player taps update again.
    private fun install(): String {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O && !packageManager.canRequestPackageInstalls()) {
            startActivity(Intent(Settings.ACTION_MANAGE_UNKNOWN_APP_SOURCES, Uri.parse("package:$packageName")))
            return "permission"
        }
        val uri = FileProvider.getUriForFile(this, "$packageName.updates", File(updatesDir(), "update.apk"))
        startActivity(
            Intent(Intent.ACTION_VIEW)
                .setDataAndType(uri, "application/vnd.android.package-archive")
                .addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION or Intent.FLAG_ACTIVITY_NEW_TASK),
        )
        return "started"
    }
}
