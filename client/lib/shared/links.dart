import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

import '../core/config.dart';
import '../l10n/app_localizations.dart';
import 'widgets.dart';

/// Server-hosted legal pages in the app language (DD-231).
String legalUrl(BuildContext context, String doc) =>
    '${AppConfig.apiBaseUrl}/$doc?lang=${Localizations.localeOf(context).languageCode}';

/// How links leave the app; tests replace it.
Future<bool> Function(Uri) launchExternal = (uri) => launchUrl(uri, mode: LaunchMode.externalApplication);

/// Opens [url] in the system browser; a toast when that is impossible.
Future<bool> openExternal(BuildContext context, String url) async {
  final ok = await launchExternal(Uri.parse(url)).catchError((_) => false);
  if (!ok && context.mounted) toast(context, Text(AppLocalizations.of(context).linkOpenFailed));
  return ok;
}
