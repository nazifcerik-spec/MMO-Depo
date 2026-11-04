import 'package:flutter_test/flutter_test.dart';
import 'package:otter/features/character/character_repository.dart';

ClassSkill _s(String code, int level, double strength, {String job = 'strike'}) =>
    ClassSkill(code: code, active: true, unlockLevel: level, job: job, strength: strength);

void main() {
  final skills = [
    _s('bolt', 1, 130, job: 'basic'),
    _s('burst', 25, 230),
    _s('element_bolt', 300, 250),
    _s('weak_bolt', 400, 200),
    _s('fire_ball', 700, 320),
    const ClassSkill(code: 'calm', active: false, unlockLevel: 50),
    const ClassSkill(code: 'old_server', active: true, unlockLevel: 1),
  ];

  List<String> kit(int level) => [for (final k in skillKit(skills, level)) k.code];

  test('the strongest unlocked skill of a job stays, stronger locked ones are ahead (DD-352)', () {
    expect(kit(25), ['bolt', 'burst', 'element_bolt', 'fire_ball', 'calm', 'old_server']);
    expect(kit(300), ['bolt', 'element_bolt', 'fire_ball', 'calm', 'old_server'], reason: 'burst is replaced');
    expect(kit(500), [
      'bolt',
      'element_bolt',
      'fire_ball',
      'calm',
      'old_server',
    ], reason: 'a weaker one replaces nothing');
    expect(kit(700), ['bolt', 'fire_ball', 'calm', 'old_server']);
  });

  test('before any skill of a job is unlocked all of them are ahead', () {
    expect(kit(1), ['bolt', 'burst', 'element_bolt', 'weak_bolt', 'fire_ball', 'calm', 'old_server']);
  });

  test('the job and its strength come from the server', () {
    final k = ClassSkill.fromJson({'code': 'a', 'kind': 'active', 'job': 'strike/x', 'strength': 250});
    expect(k.job, 'strike/x');
    expect(k.strength, 250);
    expect(ClassSkill.fromJson({'code': 'b', 'kind': 'passive'}).job, isEmpty);
  });
}
