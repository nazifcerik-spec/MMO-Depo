import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/market/market_repository.dart';

void main() {
  test('a listing reads its auction fields (DD-359)', () {
    final l = Listing.fromJson({
      'id': 'x',
      'template_code': 'iron_ore',
      'expires_at': '2030-01-01T00:00:00Z',
      'kind': 'auction',
      'buyout_copper': 900,
      'current_bid_copper': 300,
      'min_next_bid_copper': 315,
      'bid_count': 2,
      'bidder_character_id': 'me',
    });
    expect(l.isAuction, isTrue);
    expect(l.buyout, 900);
    expect(l.currentBid, 300);
    expect(l.minNextBid, 315);
    expect(l.bidCount, 2);
    expect(l.bidderId, 'me');
  });

  test('old listings are fixed price', () {
    final l = Listing.fromJson({'id': 'x', 'template_code': 'iron_ore', 'expires_at': '2030-01-01T00:00:00Z'});
    expect(l.isAuction, isFalse);
    expect(l.buyout, isNull);
  });
}
