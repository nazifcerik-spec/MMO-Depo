import 'package:flutter_test/flutter_test.dart';
import 'package:otter/shared/content_names.dart';

void main() {
  tearDown(() => ContentNames.use(const {}));

  test('names come from the table, codes are the fallback', () {
    ContentNames.use(const {'mob.grey_wolf': 'Boz Kurt'});
    expect(ContentNames.of('mob', 'grey_wolf'), 'Boz Kurt');
    expect(ContentNames.of('mob', 'cave_bat'), 'Cave Bat');
    expect(ContentNames.of('quest', 'main_01_x', 'X'), 'X');
  });

  test('affixed item names follow the language format', () {
    ContentNames.use(const {
      'item.sword_iron': 'Demir Kılıç',
      'affix.sharp': 'Keskin',
      'affix.of_the_bear': 'Ayı',
      'affix.format': '{prefix} {name} ({suffix})',
    });
    expect(ContentNames.item('sword_iron', 'sharp', 'of_the_bear'), 'Keskin Demir Kılıç (Ayı)');
    expect(ContentNames.item('sword_iron', 'sharp'), 'Keskin Demir Kılıç');
    expect(ContentNames.item('sword_iron', null, 'of_the_bear'), 'Demir Kılıç (Ayı)');
    expect(ContentNames.item('sword_iron'), 'Demir Kılıç');

    ContentNames.use(const {});
    expect(ContentNames.item('iron_ore', 'sharp', 'of_the_bear'), 'Sharp Iron Ore Of The Bear');
  });
}
