import 'package:flutter/foundation.dart';

import '../../core/api_client.dart';
import '../../core/config.dart';

const supportCategories = ['bug', 'account', 'abuse', 'feedback', 'other'];

class SupportTicket {
  const SupportTicket({
    required this.id,
    required this.category,
    required this.subject,
    required this.body,
    required this.state,
    required this.createdAt,
    this.reply,
  });

  factory SupportTicket.fromJson(Map<String, dynamic> j) => SupportTicket(
    id: j['id'] as String,
    category: j['category'] as String,
    subject: j['subject'] as String,
    body: j['body'] as String,
    state: j['state'] as String,
    reply: j['reply'] as String?,
    createdAt: DateTime.parse(j['created_at'] as String).toLocal(),
  );

  final String id;
  final String category;
  final String subject;
  final String body;
  final String state; // open | answered | closed
  final String? reply;
  final DateTime createdAt;
}

/// In-app support desk (DD-234).
class SupportRepository {
  SupportRepository(this.api);

  final ApiClient api;

  Future<List<SupportTicket>> mine() async {
    final j = await api.get('/v1/support/tickets');
    return [for (final t in (j['tickets'] as List? ?? const [])) SupportTicket.fromJson(t as Map<String, dynamic>)];
  }

  /// Attaches what helps triage (never free-form device data).
  Future<void> create({
    required String category,
    required String subject,
    required String body,
    required String locale,
  }) => api.post('/v1/support/tickets', {
    'category': category,
    'subject': subject,
    'body': body,
    'client_info': {
      'platform': kIsWeb ? 'web' : defaultTargetPlatform.name,
      'locale': locale,
      'flavor': AppConfig.flavor,
    },
  });
}
