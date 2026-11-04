// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Spanish Castilian (`es`).
class AppLocalizationsEs extends AppLocalizations {
  AppLocalizationsEs([String locale = 'es']) : super(locale);

  @override
  String get appTitle => 'MeNa World';

  @override
  String get tagline => 'Aventuras que continúan mientras no estás.';

  @override
  String get emailOrUsername => 'Correo o nombre de usuario';

  @override
  String get email => 'Correo electrónico';

  @override
  String get username => 'Nombre de usuario';

  @override
  String get password => 'Contraseña';

  @override
  String get passwordConfirm => 'Confirmar contraseña';

  @override
  String get showPassword => 'Mostrar contraseña';

  @override
  String get hidePassword => 'Ocultar contraseña';

  @override
  String get login => 'Iniciar sesión';

  @override
  String get register => 'Crear cuenta';

  @override
  String get noAccount => '¿Eres nuevo? Crea una cuenta';

  @override
  String get haveAccount => '¿Ya tienes cuenta? Inicia sesión';

  @override
  String get acceptTerms => 'Tengo 13 años o más y acepto los Términos del servicio y la Política de privacidad';

  @override
  String get fieldRequired => 'Obligatorio';

  @override
  String get invalidEmail => 'Introduce un correo válido';

  @override
  String get usernameRule => '3–32 letras, dígitos, _ o .';

  @override
  String get passwordMin => 'Al menos 8 caracteres';

  @override
  String get passwordMismatch => 'Las contraseñas no coinciden';

  @override
  String get termsRequired => 'Acepta los términos para continuar';

  @override
  String get logout => 'Cerrar sesión';

  @override
  String get retry => 'Reintentar';

  @override
  String get offline => 'No se puede conectar con el servidor. Revisa tu conexión.';

  @override
  String get characters => 'Personajes';

  @override
  String slotsUsed(int used, int cap) {
    return '$used/$cap espacios';
  }

  @override
  String get newCharacter => 'Nuevo personaje';

  @override
  String get noCharacters => 'Aún no hay personajes. Crea tu primer héroe.';

  @override
  String get create => 'Crear';

  @override
  String get characterName => 'Nombre';

  @override
  String get nameRule => 'Hasta 16 caracteres de un solo alfabeto; empieza por una letra (se admiten cifras y _)';

  @override
  String get race => 'Raza';

  @override
  String get characterClass => 'Clase';

  @override
  String get raceHuman => 'Humano';

  @override
  String get raceElf => 'Elfo';

  @override
  String get raceDwarf => 'Enano';

  @override
  String get raceOrc => 'Orco';

  @override
  String get classWarrior => 'Guerrero';

  @override
  String get classMage => 'Mago';

  @override
  String get classRogue => 'Pícaro';

  @override
  String levelLine(int level, String cls, String race) {
    return 'Nv. $level · $cls · $race';
  }

  @override
  String levelShort(int level) {
    return 'Nv. $level';
  }

  @override
  String get hp => 'PV';

  @override
  String get mp => 'PM';

  @override
  String get stamina => 'RE';

  @override
  String get energy => 'EN';

  @override
  String get xp => 'XP';

  @override
  String xpValue(int xp, int next) {
    return '$xp / $next';
  }

  @override
  String get attributes => 'Atributos';

  @override
  String get combat => 'Combate';

  @override
  String get statStr => 'FUE';

  @override
  String get statAgi => 'DES';

  @override
  String get statInt => 'INT';

  @override
  String get statVit => 'VIT';

  @override
  String get statCha => 'SAB';

  @override
  String get statLuck => 'SUE';

  @override
  String get attack => 'Ataque';

  @override
  String get defense => 'Defensa';

  @override
  String get crit => 'Crítico';

  @override
  String get hit => 'Precisión';

  @override
  String percent(int value) {
    return '$value%';
  }

  @override
  String get gold => 'Oro';

  @override
  String get switchCharacter => 'Cambiar de personaje';

  @override
  String get statePendingDelete => 'Eliminación pendiente';

  @override
  String get stateFrozen => 'Congelado';

  @override
  String get errInvalidCredentials => 'Correo, usuario o contraseña incorrectos.';

  @override
  String get errAccountLocked => 'Tu cuenta está bloqueada temporalmente.';

  @override
  String get errAccountBanned => 'Esta cuenta ha sido baneada.';

  @override
  String get errAccountExists => 'Ese correo o nombre de usuario ya está en uso.';

  @override
  String get errRateLimited => 'Demasiados intentos. Espera un momento.';

  @override
  String get errNameTaken => 'Ese nombre ya está en uso.';

  @override
  String get errSlotCap => 'Todos los espacios de personaje están ocupados.';

  @override
  String get errNotPlayable => 'Este personaje no se puede jugar ahora.';

  @override
  String get errSessionExpired => 'Tu sesión ha caducado. Inicia sesión de nuevo.';

  @override
  String errGeneric(String code) {
    return 'Algo salió mal ($code).';
  }

  @override
  String get navCharacter => 'Personaje';

  @override
  String get navActivity => 'Recolección';

  @override
  String get actMining => 'Minería';

  @override
  String get actFishing => 'Pesca';

  @override
  String get actWoodcutting => 'Tala';

  @override
  String get actHerbalism => 'Herboristería';

  @override
  String get actHunting => 'Caza';

  @override
  String get activityChoose => 'Elige una actividad';

  @override
  String get durationLabel => 'Duración';

  @override
  String get gatherResource => 'Recurso';

  @override
  String gatherXp(int xp, int next) {
    return '$xp / $next XP';
  }

  @override
  String get gatherMax => 'Nivel máximo';

  @override
  String gatherLocked(int level) {
    return 'Se desbloquea en el nivel $level';
  }

  @override
  String gatherXpGained(String name, int xp) {
    return '$name +$xp XP';
  }

  @override
  String gatherLevelUp(String name, int level) {
    return '¡$name alcanzó el nivel $level!';
  }

  @override
  String get errGatherLevel => 'Tu nivel de recolección es demasiado bajo para este recurso.';

  @override
  String startActivity(String name) {
    return 'Empezar $name';
  }

  @override
  String get collect => 'Recoger';

  @override
  String get abandon => 'Abandonar';

  @override
  String get cancel => 'Cancelar';

  @override
  String get soFar => 'Recolectado hasta ahora';

  @override
  String get activityDone => 'Completado: listo para recoger';

  @override
  String remaining(String time) {
    return 'Quedan $time';
  }

  @override
  String durHm(int h, String m) {
    return '$h h $m min';
  }

  @override
  String durM(int m) {
    return '$m min';
  }

  @override
  String durS(int s) {
    return '$s s';
  }

  @override
  String get collectEarlyTitle => '¿Recoger antes?';

  @override
  String get collectEarlyBody => 'Recibes lo recolectado hasta ahora y la actividad termina.';

  @override
  String get abandonTitle => '¿Abandonar la actividad?';

  @override
  String get abandonBody => 'Se perderá todo lo recolectado hasta ahora.';

  @override
  String get lootCollect => '¡Genial!';

  @override
  String get lootNothing => 'Nada esta vez.';

  @override
  String xpGain(int xp) {
    return '+$xp XP';
  }

  @override
  String get rarityCommon => 'Común';

  @override
  String get rarityUncommon => 'Poco común';

  @override
  String get rarityRare => 'Raro';

  @override
  String get rarityEpic => 'Épico';

  @override
  String get rarityLegendary => 'Legendario';

  @override
  String get rarityMythic => 'Mítico';

  @override
  String get errSlotBusy => 'Ya hay otra actividad en curso.';

  @override
  String get errNotActive => 'Esta actividad ya ha terminado.';

  @override
  String get evAchievement => '¡Logro desbloqueado!';

  @override
  String get evMasteryTier => '¡Nuevo nivel de maestría!';

  @override
  String get evMarketSold => 'Se ha vendido uno de tus anuncios.';

  @override
  String get evMailReceived => 'Tienes correo nuevo.';

  @override
  String get evPartyInvite => 'Te han invitado a un grupo.';

  @override
  String get evGuildInvite => 'Te han invitado a un gremio.';

  @override
  String get navBag => 'Bolsa';

  @override
  String get equipment => 'Equipo';

  @override
  String get bag => 'Bolsa';

  @override
  String weightKg(String cur, String cap) {
    return '$cur/$cap kg';
  }

  @override
  String get slotHelm => 'Casco';

  @override
  String get slotNecklace => 'Cuello';

  @override
  String get slotChest => 'Pecho';

  @override
  String get slotGloves => 'Guantes';

  @override
  String get slotLegs => 'Piernas';

  @override
  String get slotBoots => 'Botas';

  @override
  String get slotMainhand => 'Mano principal';

  @override
  String get slotOffhand => 'Mano secundaria';

  @override
  String get slotRing => 'Anillo';

  @override
  String get slotTrinket => 'Abalorio';

  @override
  String get equip => 'Equipar';

  @override
  String get unequip => 'Quitar';

  @override
  String get sell => 'Vender';

  @override
  String get sortLabel => 'Ordenar';

  @override
  String get sortSlot => 'Orden de la bolsa';

  @override
  String get sortRarity => 'Rareza';

  @override
  String get sortType => 'Tipo';

  @override
  String get sortName => 'Nombre';

  @override
  String get sortValue => 'Valor';

  @override
  String selectedCount(int count) {
    return '$count seleccionados';
  }

  @override
  String get selectAll => 'Seleccionar todo';

  @override
  String get clearSelection => 'Borrar selección';

  @override
  String sellSelectedTitle(int count) {
    return '¿Vender $count objetos?';
  }

  @override
  String destroySelectedTitle(int count) {
    return '¿Destruir $count objetos?';
  }

  @override
  String get deleteEntry => 'Eliminar';

  @override
  String craftHowMany(String name) {
    return '¿Cuántos $name?';
  }

  @override
  String craftCount(int count) {
    return 'Fabricar ×$count';
  }

  @override
  String get destroy => 'Destruir';

  @override
  String sellTitle(String name) {
    return '¿Vender $name?';
  }

  @override
  String get sellBody => 'Puedes recomprarlo al vendedor durante 2 horas.';

  @override
  String destroyTitle(String name) {
    return '¿Destruir $name?';
  }

  @override
  String get destroyBody => 'Esto no se puede deshacer.';

  @override
  String get soldFor => 'Vendido por';

  @override
  String get quantity => 'Cantidad';

  @override
  String get requiresLevel => 'Nivel requerido';

  @override
  String get weight => 'Peso';

  @override
  String weightGrams(int g) {
    return '$g g';
  }

  @override
  String get sellValue => 'Se vende por';

  @override
  String get bound => 'Ligado al alma';

  @override
  String get itemLocked => 'Este objeto está retenido por un intercambio, correo o anuncio.';

  @override
  String get errBagFull => 'Tu bolsa está llena.';

  @override
  String get errItemLocked => 'Este objeto está bloqueado.';

  @override
  String get errClassRestriction => 'Tu clase no puede usar este objeto.';

  @override
  String get errLevelRequirement => 'Tu nivel es demasiado bajo.';

  @override
  String get errNotEquippable => 'Este objeto no se puede equipar.';

  @override
  String get navAdventure => 'Aventura';

  @override
  String get segHunt => 'Caza';

  @override
  String get segDungeon => 'Mazmorra';

  @override
  String get segTravel => 'Viaje';

  @override
  String get huntChoose => 'Elige tu presa';

  @override
  String get startHunt => 'Empezar caza';

  @override
  String get huntDone => 'Caza terminada: recoge tu botín';

  @override
  String get huntEarlyBody => 'Conservas los combates ganados hasta ahora y la caza termina.';

  @override
  String mobStats(int hp, int atk, int def) {
    return 'PV $hp · ATQ $atk · DEF $def';
  }

  @override
  String huntSummary(int won, int lost) {
    return '$won ganados · $lost perdidos';
  }

  @override
  String levelUp(int level) {
    return '¡Subes de nivel! Ahora nivel $level';
  }

  @override
  String levelReq(int level) {
    return 'Nv. $level+';
  }

  @override
  String dungeonEncounters(int count, String boss) {
    return '$count combates · jefe $boss';
  }

  @override
  String get dungeonEnter => 'Entrar';

  @override
  String get dungeonDone => 'Expedición completada: reclama las recompensas';

  @override
  String get dungeonWaitHint => 'Las recompensas se desbloquean al terminar la expedición.';

  @override
  String get dungeonCleared => '¡Superada!';

  @override
  String get dungeonFailed => 'Grupo derrotado';

  @override
  String get bossKilled => 'Jefe derrotado';

  @override
  String get bonus => 'Bonificación por superarla';

  @override
  String get youAreIn => 'Estás en';

  @override
  String get routes => 'Rutas';

  @override
  String get noRoutes => 'No hay rutas desde aquí.';

  @override
  String get travelTo => 'Ir';

  @override
  String travelling(String place) {
    return 'Viajando a $place';
  }

  @override
  String get arrive => 'Llegar';

  @override
  String get arrived => 'Has llegado';

  @override
  String get cancelTravel => 'Dar la vuelta';

  @override
  String get kindCity => 'Ciudad';

  @override
  String get kindZone => 'Zona';

  @override
  String get kindWilderness => 'Tierras salvajes';

  @override
  String get kindDungeon => 'Mazmorra';

  @override
  String get kindContinent => 'Continente';

  @override
  String get safeZone => 'segura';

  @override
  String get errBusy => 'Ya hay algo en curso.';

  @override
  String get errNotReady => 'Aún no ha terminado.';

  @override
  String get errWeeklyLocked => 'Ya superaste esta mazmorra esta semana.';

  @override
  String get errNoRoute => 'No hay ruta directa hasta allí.';

  @override
  String get navWorkshop => 'Taller';

  @override
  String get segCraft => 'Fabricar';

  @override
  String get segShop => 'Tienda';

  @override
  String get profBlacksmithing => 'Herrería';

  @override
  String get profAlchemy => 'Alquimia';

  @override
  String get profCooking => 'Cocina';

  @override
  String get learn => 'Aprender';

  @override
  String recipeCount(int count) {
    return '$count recetas';
  }

  @override
  String get itemType => 'Tipo';

  @override
  String get stackSize => 'Pila';

  @override
  String get skillXpReward => 'XP de habilidad';

  @override
  String get learnHint => 'Aprende esta profesión para desbloquear sus recetas.';

  @override
  String get craft => 'Fabricar';

  @override
  String crafted(String name, int count) {
    return 'Fabricado: $name ×$count';
  }

  @override
  String skillUp(String profession, int level) {
    return '$profession ahora es nivel $level';
  }

  @override
  String requiresSkill(int level) {
    return 'Requiere habilidad $level';
  }

  @override
  String get buy => 'Comprar';

  @override
  String bought(String name) {
    return 'Comprado: $name';
  }

  @override
  String get buybacks => 'Recompra';

  @override
  String get buybackHint => 'Los objetos vendidos se pueden recomprar durante 2 horas.';

  @override
  String get errInsufficientGold => 'No tienes suficiente oro.';

  @override
  String get errInsufficientMaterials => 'Te faltan materiales.';

  @override
  String get errSkillTooLow => 'Tu habilidad es demasiado baja para esta receta.';

  @override
  String get errProfessionNotLearned => 'Aprende primero la profesión.';

  @override
  String get errPriceChanged => 'El precio ha cambiado: compruébalo e inténtalo de nuevo.';

  @override
  String get errCannotSell => 'Este objeto no se puede vender.';

  @override
  String get errBuybackExpired => 'El plazo de recompra ha terminado.';

  @override
  String get segMarket => 'Mercado';

  @override
  String get classWeapon => 'Armas';

  @override
  String get classArmor => 'Armaduras';

  @override
  String get classConsumable => 'Consumibles';

  @override
  String get classMaterial => 'Materiales';

  @override
  String get filterAll => 'Todos';

  @override
  String get sortCheapest => 'Más baratos';

  @override
  String get sortNewest => 'Más nuevos';

  @override
  String get marketEmpty => 'No hay anuncios que coincidan.';

  @override
  String listingLine(String seller, String price) {
    return '$seller · $price c/u';
  }

  @override
  String get loadMore => 'Cargar más';

  @override
  String myListings(int count) {
    return 'Mis anuncios ($count)';
  }

  @override
  String get cancelListing => 'Retirar anuncio';

  @override
  String cancelListingTitle(String name) {
    return '¿Retirar $name?';
  }

  @override
  String get cancelListingBody =>
      'El objeto vuelve a tu bolsa. Parte de la tarifa solo se reembolsa en las primeras 6 horas.';

  @override
  String buyListingTitle(String name) {
    return '¿Comprar $name?';
  }

  @override
  String buyListingBody(String price, String seller) {
    return 'Pagar $price a $seller.';
  }

  @override
  String get listOnMarket => 'Poner en el mercado';

  @override
  String listed(String name) {
    return '$name puesto en el mercado';
  }

  @override
  String quantityOf(int count, int max) {
    return 'Cantidad: $count de $max';
  }

  @override
  String get unitPrice => 'Precio por unidad';

  @override
  String suggestedPrice(String price) {
    return 'Precio de mercado: $price';
  }

  @override
  String get listingFee => 'Tarifa de anuncio (no reembolsable)';

  @override
  String get youReceive => 'Recibes al vender';

  @override
  String get errPriceOutOfRange => 'Ese precio está fuera del rango permitido.';

  @override
  String get errListingCap => 'Tienes demasiados anuncios activos.';

  @override
  String get errSelfBuy => 'No puedes comprar tu propio anuncio.';

  @override
  String get errListingUnavailable => 'Este anuncio ya no está disponible.';

  @override
  String get errNotTradable => 'Este objeto no se puede intercambiar.';

  @override
  String get social => 'Social';

  @override
  String get tabChat => 'Chat';

  @override
  String get tabFriends => 'Amigos';

  @override
  String get tabMail => 'Correo';

  @override
  String get chanGlobal => 'Mundo';

  @override
  String get chanGuild => 'Gremio';

  @override
  String get chanParty => 'Grupo';

  @override
  String get chatEmpty => 'Aún no hay mensajes. ¡Saluda!';

  @override
  String chatHint(String channel) {
    return 'Mensaje a $channel';
  }

  @override
  String get send => 'Enviar';

  @override
  String get addFriend => 'Añadir';

  @override
  String friendAdded(String name) {
    return '$name y tú ya sois amigos';
  }

  @override
  String friendRequested(String name) {
    return 'Solicitud de amistad enviada a $name';
  }

  @override
  String get friendRequests => 'Solicitudes de amistad';

  @override
  String friendsCount(int count) {
    return 'Amigos ($count)';
  }

  @override
  String get noFriends => 'Aún no hay amigos: añade a alguien por su nombre.';

  @override
  String get more => 'Más';

  @override
  String get whisper => 'Susurrar';

  @override
  String get remove => 'Quitar';

  @override
  String get block => 'Bloquear';

  @override
  String get unblock => 'Desbloquear';

  @override
  String get blockedList => 'Bloqueados';

  @override
  String removeFriendTitle(String name) {
    return '¿Quitar a $name?';
  }

  @override
  String get removeFriendBody => 'Podrás enviar otra solicitud más adelante.';

  @override
  String blockTitle(String name) {
    return '¿Bloquear a $name?';
  }

  @override
  String get blockBody =>
      'No podrá susurrarte, su chat se ocultará y su correo se descartará. También termina vuestra amistad.';

  @override
  String get accept => 'Aceptar';

  @override
  String get decline => 'Rechazar';

  @override
  String get newMail => 'Escribir correo';

  @override
  String get inboxEmpty => 'Tu buzón está vacío.';

  @override
  String get noSubject => '(sin asunto)';

  @override
  String mailFrom(String name) {
    return 'De $name';
  }

  @override
  String get attachments => 'Adjuntos';

  @override
  String get claim => 'Recoger adjuntos';

  @override
  String get mailClaimed => 'Adjuntos movidos a tu bolsa y monedero';

  @override
  String get deleteMailTitle => '¿Eliminar este correo?';

  @override
  String get deleteMailBody => 'El correo eliminado no se puede recuperar.';

  @override
  String mailSent(String postage) {
    return 'Correo enviado (franqueo $postage)';
  }

  @override
  String get mailTo => 'Para (nombre del personaje)';

  @override
  String get mailSubject => 'Asunto';

  @override
  String get mailBody => 'Mensaje';

  @override
  String get attachGold => 'Adjuntar oro';

  @override
  String get errAlreadyFriends => 'Ya sois amigos.';

  @override
  String get errAlreadyRequested => 'Ya hay una solicitud de amistad pendiente.';

  @override
  String get errBlocked => 'No puedes interactuar con este jugador.';

  @override
  String get errFriendCap => 'Tu lista de amigos está llena.';

  @override
  String get errTargetNotFound => 'No hay ningún personaje con ese nombre.';

  @override
  String get errSelfTarget => 'No puedes elegirte a ti mismo.';

  @override
  String get errInboxFull => 'Su buzón está lleno.';

  @override
  String get errTooLong => 'Es demasiado largo.';

  @override
  String get errEmpty => 'Escribe algo primero.';

  @override
  String get errNotInChannel => 'No perteneces a ese grupo.';

  @override
  String get errMuted => 'Estás silenciado y ahora no puedes chatear.';

  @override
  String get errHasAttachments => 'Recoge los adjuntos antes de eliminar.';

  @override
  String get useItem => 'Usar';

  @override
  String restoresHp(int amount) {
    return 'Restaura $amount PV';
  }

  @override
  String usedItem(String name, int hp, int max) {
    return 'Usado $name · PV $hp/$max';
  }

  @override
  String get errHpFull => 'Tu salud ya está completa.';

  @override
  String get errInCombat => 'Termina primero tu caza o expedición.';

  @override
  String get errNotUsable => 'Este objeto no se puede usar.';

  @override
  String get errTooWounded => 'Estás demasiado herido para luchar. Descansa o bebe una poción.';

  @override
  String get forgotPassword => '¿Olvidaste la contraseña?';

  @override
  String get resetHint => 'Introduce el correo de tu cuenta y te enviaremos un enlace para crear una contraseña nueva.';

  @override
  String get resetSent => 'Si existe una cuenta con ese correo, el enlace ya va de camino. Caduca en 1 hora.';

  @override
  String get sendLink => 'Enviar enlace';

  @override
  String get close => 'Cerrar';

  @override
  String get changePassword => 'Cambiar contraseña';

  @override
  String get currentPassword => 'Contraseña actual';

  @override
  String get newPassword => 'Contraseña nueva';

  @override
  String get save => 'Guardar';

  @override
  String get passwordChanged => 'Contraseña cambiada. Se cerró la sesión en tus otros dispositivos.';

  @override
  String get verifyEmailBanner => 'Confirma tu dirección de correo.';

  @override
  String get resend => 'Reenviar';

  @override
  String get verificationSent => 'Correo de verificación enviado.';

  @override
  String get errEmailNotVerified => 'Confirma tu correo antes de crear un personaje.';

  @override
  String get errAlreadyVerified => 'Tu correo ya está confirmado.';

  @override
  String get tabParty => 'Grupo';

  @override
  String get tabGuild => 'Gremio';

  @override
  String get invite => 'Invitar';

  @override
  String invitedBy(String name) {
    return 'Invitado por $name';
  }

  @override
  String get noParty => 'No estás en un grupo. Invita a alguien y se formará uno (2–5 miembros).';

  @override
  String get createParty => 'Crear grupo';

  @override
  String get pendingInvites => 'Invitaciones';

  @override
  String partyMembers(int count, int cap) {
    return 'Grupo ($count/$cap)';
  }

  @override
  String get makeLeader => 'Hacer líder';

  @override
  String get kick => 'Expulsar';

  @override
  String get leave => 'Salir';

  @override
  String get disband => 'Disolver';

  @override
  String get leavePartyTitle => '¿Salir del grupo?';

  @override
  String get leavePartyBody =>
      'Si eres el líder, le sustituye el miembro más antiguo. Un grupo con un solo miembro se disuelve.';

  @override
  String get disbandPartyTitle => '¿Disolver el grupo?';

  @override
  String get disbandPartyBody => 'Todos saldrán del grupo.';

  @override
  String get rankLeader => 'Líder';

  @override
  String get rankOfficer => 'Oficial';

  @override
  String get rankVeteran => 'Veterano';

  @override
  String get rankMember => 'Miembro';

  @override
  String get rankInitiate => 'Iniciado';

  @override
  String get createGuild => 'Fundar un gremio';

  @override
  String createGuildBody(String fee) {
    return 'Fundar un gremio cuesta $fee.';
  }

  @override
  String createGuildFee(String fee) {
    return 'Fundar gremio ($fee)';
  }

  @override
  String get guildName => 'Nombre del gremio';

  @override
  String get guildTag => 'Etiqueta';

  @override
  String get guildTagHint => '2–5 letras o dígitos';

  @override
  String guildMembers(int count, int cap) {
    return '$count/$cap miembros';
  }

  @override
  String disbandScheduled(String time) {
    return 'El gremio se disuelve en $time';
  }

  @override
  String get motd => 'Mensaje del día';

  @override
  String get editMotd => 'Editar mensaje del día';

  @override
  String get setRank => 'Cambiar rango';

  @override
  String setRankFor(String name) {
    return 'Nuevo rango para $name';
  }

  @override
  String transferTitle(String name) {
    return '¿Hacer a $name líder del gremio?';
  }

  @override
  String get transferBody => 'Pasas a ser oficial. No podrás deshacerlo tú.';

  @override
  String kickTitle(String name) {
    return '¿Expulsar a $name del gremio?';
  }

  @override
  String get leaveGuildTitle => '¿Salir del gremio?';

  @override
  String get leaveGuildBody => 'Necesitarás otra invitación para volver.';

  @override
  String get disbandGuildTitle => '¿Disolver el gremio?';

  @override
  String get disbandGuildBody => 'El gremio se disuelve en 48 horas; puedes cancelarlo hasta entonces.';

  @override
  String get cancelDisband => 'Cancelar disolución';

  @override
  String get errAlreadyInGroup => 'Ya pertenece a un grupo.';

  @override
  String get errNotLeader => 'Solo el líder puede hacer eso.';

  @override
  String get errGroupFull => 'El grupo está lleno.';

  @override
  String get errInviteInvalid => 'Esa invitación ha caducado.';

  @override
  String get errAlreadyInvited => 'Ya hay una invitación pendiente.';

  @override
  String get errGuildForbidden => 'Tu rango no puede hacer eso.';

  @override
  String get errLeaderCannotLeave => 'Cede el liderazgo o disuelve el gremio primero.';

  @override
  String get errInvalidGuildName => '3–24 letras o dígitos.';

  @override
  String get errInvalidTag => 'La etiqueta debe tener 2–5 letras o dígitos.';

  @override
  String get errTagTaken => 'Esa etiqueta ya está en uso.';

  @override
  String get progress => 'Logros y maestría';

  @override
  String get tabAchievements => 'Logros';

  @override
  String get tabMastery => 'Maestría';

  @override
  String get tabReputation => 'Reputación';

  @override
  String get tabTitles => 'Títulos';

  @override
  String achievementPoints(int points) {
    return '$points puntos de logro';
  }

  @override
  String pointsShort(int points) {
    return '$points pts';
  }

  @override
  String rewardTitle(String title) {
    return 'Título: $title';
  }

  @override
  String masteryTier(int tier, int max) {
    return 'Nivel $tier/$max';
  }

  @override
  String repDailyLeft(int amount) {
    return '+$amount hoy';
  }

  @override
  String get noTitle => 'Sin título';

  @override
  String get noTitlesYet => 'Los títulos se desbloquean con logros.';

  @override
  String get errTitleNotOwned => 'No has desbloqueado ese título.';

  @override
  String get notifications => 'Notificaciones';

  @override
  String get markAllRead => 'Marcar todo como leído';

  @override
  String get inboxNothing => 'Aún no hay notificaciones.';

  @override
  String get notifSettings => 'Ajustes de notificaciones';

  @override
  String get muteAll => 'Pausar notificaciones push';

  @override
  String get muteHint => 'Durante 24 horas. La bandeja sigue llenándose.';

  @override
  String mutedFor(String time) {
    return 'En pausa durante $time';
  }

  @override
  String get modePushInbox => 'Push + bandeja';

  @override
  String get modeInbox => 'Solo bandeja';

  @override
  String get modeSilent => 'Silencio';

  @override
  String get modeOff => 'Desactivado';

  @override
  String get nkActivity => 'Actividad terminada';

  @override
  String get nkDungeon => 'Expedición terminada';

  @override
  String get nkTravel => 'Llegada';

  @override
  String get nkWorldEvent => 'Evento mundial iniciado';

  @override
  String get nkMarketSold => 'Objeto vendido';

  @override
  String get nkMail => 'Correo nuevo';

  @override
  String get nkPartyInvite => 'Invitación de grupo';

  @override
  String get nkGuildInvite => 'Invitación de gremio';

  @override
  String get nkAchievement => 'Logro desbloqueado';

  @override
  String get nkSystem => 'Anuncio';

  @override
  String timeAgo(String time) {
    return 'Hace $time';
  }

  @override
  String get segEvents => 'Eventos';

  @override
  String get eventsEmpty => 'No hay eventos mundiales programados ahora.';

  @override
  String get worldBoss => 'Jefe mundial';

  @override
  String get invasion => 'Invasión';

  @override
  String startsIn(String time) {
    return 'Empieza en $time';
  }

  @override
  String endsIn(String time) {
    return 'Termina en $time';
  }

  @override
  String regionOnly(String region) {
    return 'Solo $region';
  }

  @override
  String get eventGoal => 'Objetivo';

  @override
  String get enlist => 'Alistarse';

  @override
  String get enlistHint =>
      'Tu poder actual se fija al alistarte. Puntúas por cada minuto del evento; las recompensas llegan al terminar.';

  @override
  String enlistedPower(int power) {
    return 'Alistado · poder $power';
  }

  @override
  String myContribution(int score, int tier) {
    return 'Tu puntuación $score · nivel de recompensa $tier';
  }

  @override
  String get eventSucceeded => 'Victoria';

  @override
  String get eventFailed => 'Fracaso';

  @override
  String tokenBalance(int tokens) {
    return '$tokens fichas de evento';
  }

  @override
  String claimDaily(int tokens) {
    return 'Reclamar $tokens';
  }

  @override
  String get claimedToday => 'Reclamado';

  @override
  String eventTokensGained(int tokens) {
    return '+$tokens fichas de evento';
  }

  @override
  String get errEventNotOpen => 'El alistamiento de este evento está cerrado.';

  @override
  String get errWrongRegion => 'Viaja primero a la región del evento.';

  @override
  String get errAlreadyEnlisted => 'Ya estás alistado.';

  @override
  String get errAlreadyClaimed => 'Ya reclamado hoy.';

  @override
  String get errEventNotActive => 'Este evento no está activo.';

  @override
  String gatherYield(String activity) {
    return 'Rendimiento de $activity';
  }

  @override
  String perHour(String value) {
    return '$value/h';
  }

  @override
  String get quests => 'Misiones';

  @override
  String get questsInProgress => 'En curso';

  @override
  String get questsNoneOpen => 'No hay misiones abiertas. Elige una abajo.';

  @override
  String get questsStory => 'Historia';

  @override
  String get questsSide => 'Misiones secundarias';

  @override
  String get questsDaily => 'Diarias';

  @override
  String resetsIn(String time) {
    return 'Se reinicia en $time';
  }

  @override
  String get turnIn => 'Entregar';

  @override
  String questRewarded(String rewards) {
    return 'Misión completada: $rewards';
  }

  @override
  String get abandonQuestTitle => '¿Abandonar la misión?';

  @override
  String get abandonQuestBody => 'Se pierde el progreso. Podrás volver a aceptarla más tarde.';

  @override
  String get abandonDailyBody => 'Se pierde el progreso y hoy no podrás volver a aceptar esta diaria.';

  @override
  String objGather(String activity) {
    return 'Recolectar ($activity)';
  }

  @override
  String get objKillAny => 'Derrotar enemigos';

  @override
  String objKill(String mob) {
    return 'Derrotar a $mob';
  }

  @override
  String get objBoss => 'Derrotar a un jefe';

  @override
  String get objDungeonAny => 'Superar una mazmorra';

  @override
  String objDungeon(String dungeon) {
    return 'Superar $dungeon';
  }

  @override
  String get objCraftAny => 'Fabricar objetos';

  @override
  String objCraft(String item) {
    return 'Fabricar $item';
  }

  @override
  String get evQuestComplete => '¡Objetivos completados: entrega la misión!';

  @override
  String get errQuestLocked => 'Termina primero la misión anterior de la historia.';

  @override
  String get errQuestTaken => 'Ya aceptaste esta misión.';

  @override
  String get errQuestNotOffered => 'Esa diaria no está en el tablón de hoy.';

  @override
  String get errQuestCap => 'Demasiadas misiones abiertas. Termina o abandona una.';

  @override
  String get errQuestNotActive => 'Esta misión no está abierta.';

  @override
  String get errQuestIncomplete => 'Los objetivos aún no están completos.';

  @override
  String get skills => 'Habilidades';

  @override
  String get skillActive => 'Activa';

  @override
  String get skillPassive => 'Pasiva';

  @override
  String skillStrike(int power) {
    return '$power% de daño';
  }

  @override
  String skillHeal(int power) {
    return 'Cura $power% de PV por debajo de la mitad';
  }

  @override
  String skillMp(int mp) {
    return '$mp PM';
  }

  @override
  String skillCooldown(int turns) {
    String _temp0 = intl.Intl.pluralLogic(turns, locale: localeName, other: 'cada $turns turnos', one: 'cada turno');
    return '$_temp0';
  }

  @override
  String get termsOfService => 'Términos del servicio';

  @override
  String get privacyPolicy => 'Política de privacidad';

  @override
  String get accountMenu => 'Cuenta';

  @override
  String get downloadMyData => 'Descargar mis datos';

  @override
  String get exportOpened => 'El enlace de descarga se abrió en tu navegador (válido 15 minutos).';

  @override
  String get linkOpenFailed => 'No se pudo abrir el enlace.';

  @override
  String get deleteAccount => 'Eliminar mi cuenta';

  @override
  String get deleteAccountBody =>
      'Tu cuenta y todos tus personajes se eliminarán definitivamente en 30 días. Iniciar sesión antes cancela la eliminación. Introduce tu contraseña para confirmar.';

  @override
  String get deleteAccountConfirm => 'Eliminar cuenta';

  @override
  String deleteScheduled(String date) {
    return 'Tu cuenta se eliminará el $date. Para cancelar, inicia sesión antes de esa fecha.';
  }

  @override
  String get deletionCancelled => 'Se canceló la solicitud de eliminación de tu cuenta.';

  @override
  String get support => 'Soporte';

  @override
  String get supportIntro =>
      '¿Has encontrado un error o necesitas ayuda con tu cuenta? Escríbenos y el equipo responderá aquí; recibirás una notificación.';

  @override
  String get supportNew => 'Nueva solicitud';

  @override
  String get supportNone => 'Aún no hay solicitudes.';

  @override
  String get supportSent => 'Solicitud recibida. Te avisaremos cuando haya respuesta.';

  @override
  String get supportCategory => 'Tema';

  @override
  String get supportSubject => 'Asunto';

  @override
  String get supportMessage => '¿Qué pasó? Paso a paso si puedes.';

  @override
  String get supportTooShort => 'El asunto necesita al menos 3 caracteres y el mensaje al menos 10.';

  @override
  String get supportBug => 'Error';

  @override
  String get supportAccount => 'Cuenta';

  @override
  String get supportAbuse => 'Abuso / denunciar a un jugador';

  @override
  String get supportFeedback => 'Sugerencia';

  @override
  String get supportOther => 'Otro';

  @override
  String get supportOpen => 'Abierta';

  @override
  String get supportAnswered => 'Respondida';

  @override
  String get supportClosed => 'Cerrada';

  @override
  String get supportReply => 'Respuesta del equipo';

  @override
  String get nkSupport => 'Respuesta de soporte';

  @override
  String get errServerBusy => 'El servidor está ocupado. Inténtalo de nuevo en unos segundos.';

  @override
  String get notifTabInbox => 'Notificaciones';

  @override
  String get notifTabRewards => 'Recompensas';

  @override
  String get rewardsNone =>
      'Aún no hay recompensas. Aquí aparece todo lo que recolectas o ganas en combate (se guarda 30 días).';

  @override
  String rewardFromActivity(String name) {
    return 'Recolección: $name';
  }

  @override
  String rewardFromHunt(String name) {
    return 'Caza: $name';
  }

  @override
  String rewardFromDungeon(String name) {
    return 'Mazmorra: $name';
  }

  @override
  String rewardLevelUp(int level) {
    return '¡Nivel $level!';
  }

  @override
  String get levelMax => 'MÁX';

  @override
  String get promotionAvailable => 'ASCENSO DISPONIBLE';

  @override
  String get promotionChoose => 'Elegir';

  @override
  String promotionConfirm(String name) {
    return '¿Convertirte en $name?';
  }

  @override
  String get promotionBody => 'Esta elección define tu clase para siempre.';

  @override
  String nextMilestone(int level, String what) {
    return 'Siguiente: Nv $level · $what';
  }

  @override
  String get msPromotion => 'Ascenso de clase';

  @override
  String msAwakening(String name) {
    return 'Despertar: $name';
  }

  @override
  String msUltimate(String name) {
    return 'Definitiva: $name';
  }

  @override
  String get msCapstone => 'Talento cúspide';

  @override
  String get msMastery => 'Maestría de clase';

  @override
  String msTitle(String name) {
    return 'Título: $name';
  }

  @override
  String estMinutes(int n) {
    return '~$n min';
  }

  @override
  String estHours(String h) {
    return '~$h h';
  }

  @override
  String overflowHeld(int xp) {
    return '$xp PX esperan tu elección';
  }

  @override
  String attributePointsBtn(int n) {
    return 'Atributos · $n puntos';
  }

  @override
  String talentPointsBtn(int n) {
    return 'Talentos · $n puntos';
  }

  @override
  String get talentsTitle => 'Talentos';

  @override
  String capstoneLocked(int level, int points) {
    return 'Requiere nivel $level y $points puntos en este árbol';
  }

  @override
  String masteryUnlocked(String title) {
    return 'Maestría de clase: $title';
  }

  @override
  String get errPromotionUnavailable => 'Ese ascenso no está disponible.';

  @override
  String get errNotEnoughPoints => 'No hay puntos suficientes.';

  @override
  String get errTalentMaxed => 'Ese talento ya está al máximo.';

  @override
  String get errCapstoneLocked => 'La cúspide requiere su nivel y un árbol completo.';

  @override
  String get errSpecRequirement => 'Tu especialización no puede usar este objeto.';

  @override
  String get errAttributeRequirement => 'Tus atributos son demasiado bajos para este objeto.';

  @override
  String get errBusyCombat => 'No mientras estés en combate, AFK o en una mazmorra.';

  @override
  String skillGuard(int pct, int turns) {
    return '+$pct% de armadura durante $turns turnos';
  }

  @override
  String skillEmpower(int pct, int turns) {
    return '+$pct% de ataque durante $turns turnos';
  }

  @override
  String get classRanger => 'Explorador';

  @override
  String get classCleric => 'Clérigo';

  @override
  String get classMonk => 'Monje';

  @override
  String get classPaladin => 'Paladín';

  @override
  String get classDruid => 'Druida';

  @override
  String get classBard => 'Bardo';

  @override
  String get classShaman => 'Chamán';

  @override
  String get errEquipmentPermission => 'Tu clase no puede usar este tipo de equipo.';

  @override
  String get classResource => 'Recurso';

  @override
  String get classFocus => 'Atributo principal';

  @override
  String classSecondary(String names) {
    return 'Secundario: $names';
  }

  @override
  String get resStartsEmpty => 'Empieza vacío';

  @override
  String get resStartsFull => 'Empieza lleno';

  @override
  String resStartsAt(int pct) {
    return 'Empieza al $pct%';
  }

  @override
  String resPerHit(int n) {
    return '+$n por golpe';
  }

  @override
  String resPerTurn(int n) {
    return '+$n por turno';
  }

  @override
  String baseAttributesHint(int total) {
    return 'Atributos base · 1 punto cada 2 niveles, $total en total';
  }

  @override
  String get classEquipment => 'Equipo';

  @override
  String get firstSkills => 'Primeras habilidades';

  @override
  String get classPath => 'Ruta de clase';

  @override
  String get classPathHint => 'Toca una clase para ver sus habilidades.';

  @override
  String classMilestones(int awakening, int ultimate, int capstone, int mastery) {
    return 'Nv $awakening Despertar · Nv $ultimate Definitiva · Nv $capstone talento cúspide · Nv $mastery maestría de clase';
  }

  @override
  String get skillTagActive => 'activa';

  @override
  String get skillTagPassive => 'pasiva';

  @override
  String get skillTagAwakening => 'despertar';

  @override
  String get skillTagUltimate => 'definitiva';

  @override
  String skillCostCd(int cost, int cd) {
    return 'Coste $cost · reutilización $cd turnos';
  }

  @override
  String masteryAt(int level, String title) {
    return 'Maestría Nv $level: $title';
  }

  @override
  String get classInfoFailed => 'No se pudieron cargar los detalles de la clase.';

  @override
  String xpShare(int pct) {
    return 'PX $pct%';
  }

  @override
  String get dungeonsSolo => 'Mazmorras';

  @override
  String get raidsTab => 'Incursiones';

  @override
  String raidGroup(int min, int max) {
    return 'Grupo de $min a $max';
  }

  @override
  String raidBossGear(int count) {
    return 'Jefe: $count piezas de equipo garantizadas';
  }

  @override
  String get raidOpen => 'Abrir incursión';

  @override
  String get raidNeedParty =>
      'Las incursiones son para un grupo. Ponte en cola para uno al azar, únete a un anuncio o publica el tuyo.';

  @override
  String get raidLeaderOnly => 'El líder del grupo abre las incursiones.';

  @override
  String get raidLobby => 'Sala de la incursión';

  @override
  String get raidJoin => 'Unirse';

  @override
  String get raidLeave => 'Salir';

  @override
  String get raidCancel => 'Cancelar';

  @override
  String get raidCancelBody => 'La sala se cierra para todos.';

  @override
  String get raidStart => 'Iniciar incursión';

  @override
  String get raidStartBody => 'Una vez iniciada nadie puede unirse ni salir.';

  @override
  String raidMembers(int count, int max) {
    return 'Miembros ($count/$max)';
  }

  @override
  String raidNeedMore(int min) {
    return 'Se necesitan al menos $min miembros';
  }

  @override
  String get raidWaitHint => 'Las recompensas se abren al terminar; cada miembro recoge las suyas.';

  @override
  String get raidDone => 'Incursión terminada: recoge tus recompensas';

  @override
  String get raidLocked => 'Completada esta semana';

  @override
  String get raidCollected => 'recogido';

  @override
  String get errRaidNeedsParty => 'Únete primero a un grupo.';

  @override
  String get errRaidExists => 'Tu grupo ya tiene una incursión abierta.';

  @override
  String get errAlreadyInRaid => 'Ya estás en una incursión.';

  @override
  String get errRaidNotOpen => 'Esta incursión ya no está abierta.';

  @override
  String get errRaidNotEnough => 'No hay suficientes miembros para empezar.';

  @override
  String get errRaidWeeklyLocked => 'Ya completaste esta incursión esta semana.';

  @override
  String errCharacterBusy(String action) {
    return 'Termina primero lo que hace tu personaje: $action.';
  }

  @override
  String get busyGathering => 'recolección';

  @override
  String get busyHunt => 'caza';

  @override
  String get busyDungeon => 'mazmorra';

  @override
  String get busyRaid => 'incursión';

  @override
  String get busyTravel => 'viaje';

  @override
  String get busyCrafting => 'fabricación';

  @override
  String get errCraftActive => 'Ya hay una fabricación en curso.';

  @override
  String get errCraftTooLong => 'Ese trabajo duraría más de 6 horas.';

  @override
  String craftEach(String time) {
    return '$time cada uno';
  }

  @override
  String craftTotal(String time) {
    return 'Total $time';
  }

  @override
  String craftRunning(String item, int count) {
    return 'Fabricando $item ×$count';
  }

  @override
  String craftProgress(int done, int count) {
    return '$done/$count listos';
  }

  @override
  String get craftDone => 'Todo listo: recoge';

  @override
  String get craftCollectEarlyTitle => '¿Recoger ahora?';

  @override
  String craftCollectEarlyBody(int done, int rest) {
    return 'Se entregan los $done objetos terminados; los materiales y el oro de los otros $rest se devuelven.';
  }

  @override
  String craftCollected(int made, int refunded) {
    return '$made fabricados, $refunded devueltos';
  }

  @override
  String get craftBusyHint => 'Fabricar lleva tiempo; mientras tanto tu personaje no hace nada más.';

  @override
  String get eventHowTitle => 'Cómo participar';

  @override
  String get eventRewardsTitle => 'Recompensas';

  @override
  String eventStepEnlist(int level) {
    return 'Inscríbete desde 24 horas antes de que empiece hasta que termine (nivel $level+).';
  }

  @override
  String eventStepRegion(String region) {
    return 'Debes estar en $region al inscribirte.';
  }

  @override
  String eventStepPower(int power) {
    return 'Tu poder se fija al inscribirte: ataque × 2 + defensa + PV ÷ 10. El tuyo ahora: $power.';
  }

  @override
  String eventStepScore(int tier) {
    return 'Puntos = poder × minutos de asistencia, desde el inicio (o desde que te inscribes, si es después). Se detiene en el nivel $tier.';
  }

  @override
  String get eventStepGoal =>
      'Si entre todos se alcanza el objetivo, las recompensas llegan al nivel máximo; si no, se detienen en el nivel 2.';

  @override
  String get eventStepPaid => 'Las recompensas se pagan solas al cerrar. Inscribirte no detiene tus otras acciones.';

  @override
  String eventTier(int tier, String score) {
    return 'Nivel $tier · $score+ puntos';
  }

  @override
  String eventTierTokens(int count) {
    return '+$count fichas de evento';
  }

  @override
  String eventEstimate(String score, int tier) {
    return 'Si te unes ahora: unos $score puntos → nivel $tier';
  }

  @override
  String eventEstimateNone(String score) {
    return 'Si te unes ahora: unos $score puntos, no alcanza un nivel';
  }

  @override
  String seasonalHow(int count) {
    return 'Mientras dure el evento, reclama $count fichas una vez al día (el día empieza a las 00:00 UTC).';
  }

  @override
  String seasonalBenefit(int count) {
    return 'Las fichas se guardan en tu cuenta (saldo: $count).';
  }

  @override
  String get eventDetails => 'Detalles';

  @override
  String get createCharacterTitle => 'Crear personaje';

  @override
  String get classTraits => 'Rasgos de la clase';

  @override
  String classResourceLine(String name) {
    return 'Recurso: $name';
  }

  @override
  String classFocusLine(String primary, String secondary) {
    return 'Atributo principal: $primary · segundo: $secondary';
  }

  @override
  String classFocusOnly(String primary) {
    return 'Atributo principal: $primary';
  }

  @override
  String get classDifficulty => 'Dificultad';

  @override
  String get difficultyEasy => 'Fácil — indulgente, segura en cacerías AFK';

  @override
  String get difficultyMedium => 'Media — equilibrada; importa el orden de habilidades';

  @override
  String get difficultyHard => 'Difícil — frágil; necesita el equipo adecuado';

  @override
  String get skillsTapHint => 'toca para ver detalles';

  @override
  String get chooseRace => 'Elige tu raza';

  @override
  String continueAs(String name) {
    return 'Continuar como $name';
  }

  @override
  String get nameYourCharacter => 'Ponle nombre a tu personaje';

  @override
  String get raceTroll => 'Trol';

  @override
  String get raceUndead => 'No muerto';

  @override
  String get raceDarkElf => 'Elfo oscuro';

  @override
  String get raceDrakari => 'Drakari';

  @override
  String get raceTauren => 'Tauren';

  @override
  String get raceGnome => 'Gnomo';

  @override
  String foodBuffLine(int atk, int def, int minutes) {
    return '+$atk% ataque y +$def% defensa durante $minutes min';
  }

  @override
  String fedToast(int atk, int def) {
    return 'Saciado: +$atk% ataque, +$def% defensa';
  }

  @override
  String fedActive(int atk, int def, String time) {
    return 'Saciado · +$atk% ataque · +$def% defensa · quedan $time';
  }

  @override
  String huntReach(String regions) {
    return 'Zonas de caza desde aquí: $regions';
  }

  @override
  String get huntNoMobsHere => 'Aquí no hay nada que cazar. Viaja primero a una zona de caza.';

  @override
  String get huntOtherRegions => 'Monstruos de otras zonas';

  @override
  String get huntTravelHint => 'Para cazar un monstruo, viaja a su zona o a una zona unida a ella por camino.';

  @override
  String get errMobRegion => 'Este monstruo vive en otra zona. Viaja allí primero.';

  @override
  String restoresHpPct(int pct) {
    return 'Restaura el $pct% de la vida máxima';
  }

  @override
  String get statEvasion => '% Evasión';

  @override
  String get statBlock => '% Bloqueo';

  @override
  String get statCritDmg => '% Daño crítico';

  @override
  String get statLifesteal => '% Robo de vida';

  @override
  String get statDmgReduction => '% Reducción de daño';

  @override
  String get classJewelry => 'Joyería';

  @override
  String get mapHint =>
      'Puedes viajar directamente a cualquier zona que tu nivel permita; el tiempo sigue el camino más corto. Arrastra y amplía el mapa; toca una zona.';

  @override
  String get mapNoRoad => 'Ningún camino directo lleva aquí. Viaja por las zonas que brillan a tu lado.';

  @override
  String get youAreHere => 'Estás aquí.';

  @override
  String travelTime(String time) {
    return 'Tiempo de viaje: $time';
  }

  @override
  String xpPerKill(int xp) {
    return '$xp XP';
  }

  @override
  String get deleteCharacter => 'Eliminar personaje';

  @override
  String deleteCharacterBody(String name) {
    return '$name se eliminará ahora mismo, con todos sus objetos y oro. No se puede deshacer. Escribe el nombre del personaje para confirmar.';
  }

  @override
  String get characterNameConfirm => 'Nombre del personaje';

  @override
  String characterDeletesAt(String when) {
    return 'Se elimina el $when';
  }

  @override
  String get cancelCharacterDelete => 'Cancelar eliminación';

  @override
  String characterDeleted(String name) {
    return '$name fue eliminado.';
  }

  @override
  String characterDeleteCancelled(String name) {
    return '$name ya no se eliminará.';
  }

  @override
  String get errNameMismatch => 'El nombre no coincide.';

  @override
  String attrGainPerPoint(String gains) {
    return 'Por punto: $gains';
  }

  @override
  String attrGainPending(String gains) {
    return 'Se añade: $gains';
  }

  @override
  String get attrNoEffect => 'Sin efecto para tu clase';

  @override
  String get attrInputHint => 'Escribe cuántos puntos recibe cada atributo';

  @override
  String get gainHp => 'Vida';

  @override
  String get gainMp => 'Maná';

  @override
  String get gainHpRegen => 'Vida/s';

  @override
  String get huntSupplies => 'Suministros';

  @override
  String get huntSuppliesHint =>
      'Las pociones se beben entre combates con menos de la mitad de vida, la comida cuando acaba su bonificación y las de maná con poco maná. Lo que sobre vuelve a tu bolsa.';

  @override
  String get supplyHp => 'Pociones de vida';

  @override
  String get supplyMp => 'Pociones de maná';

  @override
  String get supplyMpHint => 'Para clases que usan maná, energía, concentración, energía natural o espíritu';

  @override
  String get supplyFood => 'Comida';

  @override
  String get huntNoSupplies => 'No tienes pociones ni comida: cómpralas en la tienda o fabrícalas.';

  @override
  String restoresMpPct(int pct) {
    return 'Rellena el $pct% de tu maná';
  }

  @override
  String get huntSupplyOnly => 'Solo se usa como suministro de caza';

  @override
  String get errBadSupply => 'Ese objeto no se puede llevar de caza.';

  @override
  String get errNotEnoughSupplies => 'No tienes suficientes de ese suministro.';

  @override
  String get errHuntOnly => 'Llévalo de caza como suministro.';

  @override
  String get huntFled => 'Tu vida bajó del 5 %: abandonaste el combate';

  @override
  String suppliesUsed(String list) {
    return 'Usado: $list';
  }

  @override
  String get compareTitle => 'Comparado con lo que llevas';

  @override
  String compareWorn(String name) {
    return 'Equipado: $name';
  }

  @override
  String get compareEmpty => 'Esa ranura está vacía: al equiparlo se suman todos.';

  @override
  String get compareColStat => 'Atributo';

  @override
  String get compareColWorn => 'Equipado';

  @override
  String get compareColThis => 'Este';

  @override
  String get compareColDiff => 'Cambio';

  @override
  String get compareSame => 'Mismos atributos que lo que llevas.';

  @override
  String get wearOk => 'Tu clase puede llevarlo';

  @override
  String get wearClass => 'Tu clase no puede llevarlo';

  @override
  String wearLevel(int level) {
    return 'Requiere nivel $level';
  }

  @override
  String get wearAttribute => 'Tus atributos no bastan';

  @override
  String get wearLegendClass => 'no es para tu clase';

  @override
  String get wearLegendLater => 'aún no (nivel o atributos)';

  @override
  String get enchantAction => 'Encantar';

  @override
  String get enchantPick => 'Elige el equipo a encantar';

  @override
  String get enchantNone => 'Ningún equipo encaja con este pergamino.';

  @override
  String enchantFits(String groups) {
    return 'Se aplica a: $groups';
  }

  @override
  String get enchantGroupWeapon => 'armas';

  @override
  String get enchantGroupArmor => 'armadura';

  @override
  String get enchantGroupJewelry => 'joyas';

  @override
  String get enchantGroupAny => 'cualquier equipo';

  @override
  String get enchantReplaceTitle => '¿Reemplazar el encantamiento?';

  @override
  String enchantReplaceBody(String old, String next) {
    return '$old se pierde y $next ocupa su lugar.';
  }

  @override
  String enchantDone(String item) {
    return '$item encantado';
  }

  @override
  String enchantOn(String name) {
    return 'Encantamiento: $name';
  }

  @override
  String get wornTag => 'equipado';

  @override
  String get errNotEnchant => 'No es un pergamino de encantamiento.';

  @override
  String get errEnchantSlot => 'El encantamiento no encaja en ese equipo.';

  @override
  String get statGatherMining => 'Rendimiento minero';

  @override
  String get statGatherHerbalism => 'Rendimiento de herboristería';

  @override
  String get statGatherWoodcutting => 'Rendimiento de tala';

  @override
  String get statGatherFishing => 'Rendimiento de pesca';

  @override
  String get statGatherHunting => 'Rendimiento de caza';

  @override
  String get statCraftSpeed => 'Velocidad de fabricación';

  @override
  String diffXp(int pct) {
    return '$pct % XP';
  }

  @override
  String get diffNoXp => 'sin XP';

  @override
  String get diffLegend => 'Color de receta: la XP que da a tu nivel';

  @override
  String get recipeGroupEnchantWeapon => 'Encantamientos de arma';

  @override
  String get recipeGroupEnchantArmor => 'Encantamientos de armadura';

  @override
  String get recipeGroupEnchantJewelry => 'Encantamientos de joyería';

  @override
  String get recipeGroupEnchantProfession => 'Encantamientos de oficio';

  @override
  String get partyHint =>
      'Los miembros que hacen el mismo trabajo a la vez (la misma recolección, profesión, caza o mazmorra) ganan cada uno +3 % PX por miembro: cinco juntos, +15 %. Nadie cede PX.';

  @override
  String get partyChat => 'Chat de grupo';

  @override
  String partyShared(String names, int kept, int earned) {
    return 'Bonificación de grupo con $names: $earned → $kept PX';
  }

  @override
  String evRaidMatched(String raid) {
    return 'Se encontró un grupo para $raid. Mira Incursiones.';
  }

  @override
  String get nkRaidMatched => 'Grupo de incursión';

  @override
  String get workHunt => 'Caza';

  @override
  String get workDungeon => 'Mazmorra';

  @override
  String get raidFind => 'Buscar grupo';

  @override
  String get raidQueueJoin => 'Hacer cola';

  @override
  String get raidQueueLeave => 'Salir de la cola';

  @override
  String raidQueued(String raid, int waiting) {
    return 'En cola para $raid · $waiting esperando';
  }

  @override
  String get raidQueueHint => 'Se forma un grupo al azar cuando esperan 5, o tras 5 minutos con al menos 2.';

  @override
  String get raidAdvertise => 'Anunciar';

  @override
  String get raidUnlist => 'Retirar anuncio';

  @override
  String raidListed(String raid) {
    return 'Tu grupo está en el tablón: $raid';
  }

  @override
  String get raidBoard => 'Tablón de anuncios';

  @override
  String get raidBoardEmpty => 'No hay anuncios abiertos.';

  @override
  String get raidJoinParty => 'Unirse';

  @override
  String get raidPickRaid => 'Incursión';

  @override
  String get errListingGone => 'Ese anuncio ya no existe o está lleno.';

  @override
  String get workGathering => 'Recolección';

  @override
  String get workAdventure => 'Aventura';

  @override
  String get workRaid => 'Incursión';

  @override
  String get workTravel => 'Viajando';

  @override
  String get workCrafting => 'Fabricación';

  @override
  String get workIdle => 'Inactivo';

  @override
  String get errDungeonRegion => 'Primero viaja a este lugar en el mapa.';

  @override
  String get dungeonsElsewhere =>
      'Más mazmorras e incursiones se abren al viajar: cada una se entra desde su región del mapa.';

  @override
  String get dungeonNoneHere => 'No hay mazmorras aquí. Viaja por el mapa a una región con mazmorra.';

  @override
  String dungeonXp(int xp) {
    return '≈ $xp PX para ti';
  }

  @override
  String appVersionLine(String version) {
    return 'Versión $version';
  }

  @override
  String get updateButton => 'Actualizar';

  @override
  String updateAvailable(String version) {
    return 'Nueva versión lista: $version';
  }

  @override
  String get updateLatest => 'Tienes la última versión';

  @override
  String updateConfirm(String version, int size) {
    return '¿Descargar e instalar la versión $version ($size MB)?';
  }

  @override
  String get updateDownloading => 'Descargando…';

  @override
  String get updatePermission => 'Permite instalaciones desde MeNa World y vuelve a tocar Actualizar.';

  @override
  String get updateFailed => 'No se pudo obtener la actualización. Revisa tu conexión.';

  @override
  String get fightReport => 'Informe de combate';

  @override
  String get endTimeUp => 'Se acabó el tiempo: duró hasta el final previsto.';

  @override
  String get endEarly => 'Recogido antes: se detuvo antes del final previsto.';

  @override
  String endFled(String mob) {
    return 'Tu vida bajó del 5 % contra $mob: te retiraste.';
  }

  @override
  String endDefeat(String mob) {
    return 'Caíste ante $mob.';
  }

  @override
  String endPartyDefeat(String mob) {
    return 'El grupo cayó ante $mob.';
  }

  @override
  String endTimeout(String mob) {
    return 'El combate con $mob superó el límite de turnos: se detuvo.';
  }

  @override
  String get endCleared => 'Superada: cayeron todas las oleadas y el jefe.';

  @override
  String get repCampHint =>
      'Por debajo del 50 % de vida el cazador acampa hasta llenarla antes del siguiente combate, así que rara vez baja al 5 %.';

  @override
  String get repYou => 'Tú';

  @override
  String get repLevel => 'Nivel';

  @override
  String get repHp => 'Vida';

  @override
  String get repAttack => 'Ataque';

  @override
  String get repDefense => 'Defensa';

  @override
  String get repEvasion => 'Evasión';

  @override
  String get repHitChance => 'Prob. de acierto';

  @override
  String get repCrit => 'Crítico';

  @override
  String repFights(int n) {
    return '$n combates';
  }

  @override
  String repSwings(int hits, int swings, int pct, int avg, int max) {
    return '$hits/$swings aciertos ($pct %), media $avg, máx. $max';
  }

  @override
  String get repYourSwings => 'Tus golpes';

  @override
  String get repPartySwings => 'Golpes del grupo';

  @override
  String get repEnemySwings => 'Golpes del enemigo';

  @override
  String repLowest(int pct) {
    return 'Vida más baja en combate: $pct %';
  }

  @override
  String repRest(String time) {
    return 'Acampado: $time';
  }

  @override
  String get restart => 'Reiniciar';

  @override
  String get restartTitle => '¿Recoger y empezar de nuevo?';

  @override
  String get restartBody => 'Lo obtenido hasta ahora va a tu bolsa y lo mismo vuelve a empezar con la misma duración.';

  @override
  String get dailyTitle => 'Recompensa diaria';

  @override
  String get dailyBody =>
      'Reclama una vez al día con cualquier personaje. Cada semana seguida da más y el día 29 lo máximo; si faltas un día, el ciclo vuelve a empezar.';

  @override
  String get dailyClaim => 'Reclamar';

  @override
  String get dailyTomorrow => 'Vuelve mañana';

  @override
  String dailyDay(int n) {
    return 'Día $n';
  }

  @override
  String get dailyOpen => 'Recompensa diaria';

  @override
  String chatMutedNotice(String name, int min) {
    return '$name fue silenciado $min minutos por insultar.';
  }

  @override
  String get storyTitle => 'Historia';

  @override
  String get storyOpen => 'Modo historia';

  @override
  String get storyComingSoon =>
      'El modo historia llegará pronto. Los capítulos aparecerán aquí a medida que se escriban.';

  @override
  String storyLocked(int level) {
    return 'Se abre en el nivel $level';
  }

  @override
  String get bankTitle => 'Banco';

  @override
  String get bankHint =>
      'Compartido por todos tus personajes. Toca un objeto de tu bolsa para depositarlo o uno del banco para sacarlo. Los objetos ligados se quedan con su personaje.';

  @override
  String get bankEmpty => 'El banco está vacío.';

  @override
  String get bankDeposit => 'Depositar';

  @override
  String get bankWithdraw => 'Sacar';

  @override
  String get errBankFull => 'El banco está lleno.';

  @override
  String get tabTrade => 'Intercambio';

  @override
  String get tradeHint =>
      'Intercambia objetos con otro jugador. Ambos ponen objetos y confirman; cada lado paga un impuesto del 8 % del precio medio de mercado de lo que da, por unidad.';

  @override
  String get tradeWith => 'Nombre del jugador';

  @override
  String get tradeOpen => 'Iniciar intercambio';

  @override
  String tradeTitle(String name) {
    return 'Intercambio con $name';
  }

  @override
  String get tradeYours => 'Tú das';

  @override
  String tradeTheirs(String name) {
    return '$name da';
  }

  @override
  String get tradeNothing => 'Nada todavía.';

  @override
  String get tradeAdd => 'Añadir objeto';

  @override
  String get tradeRemove => 'Retirar';

  @override
  String get tradeTax => 'Tu impuesto';

  @override
  String get tradeReady => 'Listo';

  @override
  String get tradeWaiting => 'No listo';

  @override
  String get tradeConfirm => 'Confirmar intercambio';

  @override
  String get tradeWaitingThem => 'Esperando al otro jugador';

  @override
  String get tradeCancel => 'Cancelar intercambio';

  @override
  String get tradeCancelBody => 'No se intercambia nada.';

  @override
  String get tradeDone => 'Intercambio completado.';

  @override
  String get errTradeBusy => 'Uno de vosotros ya está intercambiando.';

  @override
  String get errTradeClosed => 'Ese intercambio terminó.';

  @override
  String get errTradeItem => 'Ese objeto no se puede intercambiar ahora.';

  @override
  String get errTradeEmpty => 'Primero pon algo.';

  @override
  String get errTradeTooMany => 'Demasiados objetos en el intercambio.';

  @override
  String get errCraftLevel => 'El nivel de tu personaje es inferior al del objeto.';

  @override
  String get gatherNoneHere => 'Tus recursos crecen en otros mapas: viaja a la región indicada.';

  @override
  String get errGatherRegion => 'Ese recurso crece en otro mapa.';

  @override
  String get helpersOpen => 'Ayudantes';

  @override
  String get helpersTitle => 'Ayudantes';

  @override
  String get helpersHint =>
      'Contrata ayudantes que recolecten y fabriquen a tu lado. Trabajan a la mitad de tu velocidad y ganan la mitad de la PX, que va a tu profesión; cada tipo tiene su bonificación.';

  @override
  String get helpersGather => 'Recolección';

  @override
  String get helpersCraft => 'Fabricación';

  @override
  String get helperHire => 'Contratar';

  @override
  String helperHireBody(String name, String price) {
    return '¿Contratar al $name por $price?';
  }

  @override
  String get helperSend => 'Enviar a trabajar';

  @override
  String get helperCancel => 'Llamar de vuelta';

  @override
  String get helperCancelBody => 'Una recolección no trae nada; una fabricación devuelve materiales y oro.';

  @override
  String get helperDone => 'Hecho: recógelo';

  @override
  String helperWorking(String time) {
    return 'Trabajando, faltan $time';
  }

  @override
  String helperCollected(int items, int xp, String profession) {
    return '$items objetos, +$xp PX de $profession';
  }

  @override
  String get helperNoRecipes => 'Aún no hay recetas a tu nivel en esta profesión.';

  @override
  String get helperPickRecipe => 'Elige una receta';

  @override
  String helperBonusYield(int pct) {
    return '+$pct % de rendimiento';
  }

  @override
  String helperBonusSpeed(int pct) {
    return '$pct % más rápido';
  }

  @override
  String helperBonusXp(int pct) {
    return '+$pct % PX';
  }

  @override
  String helperBonusGold(int pct) {
    return '$pct % menos de coste en oro';
  }

  @override
  String helperBonusDouble(int pct) {
    return '$pct % de probabilidad de doble producción';
  }

  @override
  String get helperMiner => 'Minero';

  @override
  String get helperWoodsman => 'Leñador';

  @override
  String get helperAngler => 'Pescador';

  @override
  String get helperHerbalist => 'Herbolario';

  @override
  String get helperTrapper => 'Trampero';

  @override
  String get helperSmith => 'Herrero';

  @override
  String get helperTanner => 'Curtidor';

  @override
  String get helperWeaver => 'Tejedor';

  @override
  String get helperCarver => 'Tallador';

  @override
  String get helperApothecary => 'Boticario';

  @override
  String get helperCook => 'Cocinero';

  @override
  String get helperSmelter => 'Fundidor';

  @override
  String get helperScribe => 'Escriba';

  @override
  String get helperGemcutter => 'Tallador de gemas';

  @override
  String get errHelperBusy => 'Ese ayudante ya está trabajando.';

  @override
  String get errHelperNotReady => 'El ayudante aún no ha terminado.';

  @override
  String get shopOpen => 'Tienda';

  @override
  String get shopTitle => 'Tienda de créditos';

  @override
  String get shopCredits => 'Créditos';

  @override
  String get shopAvatars => 'Avatares';

  @override
  String get shopBackgrounds => 'Fondos';

  @override
  String shopBalance(int n) {
    return '$n créditos';
  }

  @override
  String get shopBuyCredits => 'Compra créditos para potenciadores y cosméticos.';

  @override
  String get shopPaymentsSoon => 'La compra de créditos con dinero llegará pronto.';

  @override
  String get shopBoosts => 'Potenciadores de PX';

  @override
  String get shopBuy => 'Comprar';

  @override
  String shopBuyBody(int n) {
    return '¿Gastar $n créditos?';
  }

  @override
  String shopPrice(int n) {
    return '$n cr';
  }

  @override
  String get boostGather => 'Recolección';

  @override
  String get boostAdventure => 'Combate';

  @override
  String get boostCraft => 'Profesiones';

  @override
  String get boostDesc => '+100 % PX durante 2 horas';

  @override
  String get boostShort => '+100 % PX, 2 h';

  @override
  String boostLeft(String time) {
    return 'Activo, faltan $time';
  }

  @override
  String boostWatchAd(int n) {
    return 'Ver anuncio (quedan $n)';
  }

  @override
  String get boostAdsSoon => 'Anuncios pronto';

  @override
  String get errInsufficientCredits => 'Créditos insuficientes.';

  @override
  String get errBoostActive => 'Ese potenciador ya está activo.';

  @override
  String get errAdLimit => 'No quedan anuncios hoy.';

  @override
  String helperCanMake(int n) {
    return 'Puede fabricar: $n';
  }

  @override
  String get helperPickResource => 'Qué recolectar';

  @override
  String get helperCraftHint => 'Los materiales y el oro se cobran ahora; mientras tanto sigues jugando.';

  @override
  String get petsOpen => 'Mascotas';

  @override
  String get petsHint =>
      'Los huevos de mascota caen al superar bandas de nivel 400 o más. El huevo eclosiona en la incubadora; la mascota sube de nivel con sus propias cacerías y tu mascota activa te da su bonificación. Mientras está activa también lucha a tu lado y añade daño a cada golpe.';

  @override
  String get petsEggs => 'Huevos';

  @override
  String get petsNoEggs => 'No tienes huevos. Supera bandas de nivel 400 o más.';

  @override
  String petsMine(int n, int max) {
    return 'Mis mascotas ($n/$max)';
  }

  @override
  String get petHatch => 'Incubar';

  @override
  String petHatchBody(String egg, String time) {
    return '¿Incubar un $egg? Eclosiona en $time.';
  }

  @override
  String get petHatching => 'En su huevo';

  @override
  String get petHatched => '¡Ha eclosionado!';

  @override
  String petPower(int n) {
    return 'Lucha como un monstruo de nivel $n';
  }

  @override
  String get petActive => 'Activa';

  @override
  String get petActivate => 'Activar';

  @override
  String get petDeactivate => 'Descansar';

  @override
  String get petHunt => 'Enviar a cazar';

  @override
  String petHunting(String mob) {
    return 'Cazando $mob';
  }

  @override
  String get petHuntDone => 'Cacería terminada';

  @override
  String get petCollectEarly => 'Terminar ahora paga los combates hasta el momento.';

  @override
  String petCollected(int wins, int xp) {
    return '$wins victorias, +$xp PX de mascota';
  }

  @override
  String petLevelUp(String name, int level) {
    return '¡$name alcanzó el nivel $level!';
  }

  @override
  String get petDied => 'Tu mascota fue derrotada; la cacería terminó antes.';

  @override
  String get petPickMob => 'Monstruo a cazar';

  @override
  String petXpPerWin(int n) {
    return '+$n PX por victoria';
  }

  @override
  String get errPetCap => 'Ya tienes el máximo de mascotas.';

  @override
  String get errPetHatching => 'La mascota aún no ha eclosionado.';

  @override
  String get errPetHunting => 'La mascota ya está cazando.';

  @override
  String petStageLevels(int from, int to) {
    return 'Nv $from–$to';
  }

  @override
  String get regionMobs => 'Monstruos';

  @override
  String get regionDungeons => 'Mazmorras';

  @override
  String get regionGather => 'Recursos';

  @override
  String get regionEmpty => 'Ningún monstruo, mazmorra, recurso ni misión pertenece a esta región.';

  @override
  String get regionRaids => 'Incursiones';

  @override
  String get tagUse => 'Uso';

  @override
  String get tagAllClasses => 'Todas las clases';

  @override
  String get tagDrops => 'Lo sueltan';

  @override
  String get tagRaid => 'Incursión';

  @override
  String get tagDungeon => 'Mazmorra';

  @override
  String get tagCrafted => 'Lo fabrica';

  @override
  String tagSkill(int n) {
    return 'habilidad $n';
  }

  @override
  String get tagSold => 'Lo vende';

  @override
  String get tagGathered => 'Se recolecta con';

  @override
  String get tagUsedIn => 'Material para';

  @override
  String get petEvolution => 'Evolución';

  @override
  String get petStrength => 'Detalles de fuerza';

  @override
  String petAssistRange(int min, int max) {
    return 'Apoyo en combate: +$min% – +$max% de daño en cada golpe';
  }

  @override
  String petBonusRange(String min, String max) {
    return 'Bonificación: $min – $max';
  }

  @override
  String petFightsRange(int min, int max) {
    return 'Caza como un monstruo de nivel $min – $max';
  }

  @override
  String get petFind => 'Dónde conseguirla';

  @override
  String petAssist(int pct) {
    return 'Apoyo en combate: +$pct% de daño';
  }

  @override
  String get petLevelNote => 'Del nivel 1 al 100';

  @override
  String get petIndex => 'Índice de mascotas';

  @override
  String petIndexCount(int n, int max) {
    return '$n / $max descubiertas';
  }

  @override
  String slotCount(int n, int max) {
    return '$n / $max';
  }

  @override
  String get petEggTapHint => 'Toca un huevo para incubarlo.';

  @override
  String get petTapHint => 'Toca una mascota para verla · el punto verde marca la activa.';

  @override
  String get petIndexHint => 'Toca una mascota para abrir su página.';

  @override
  String get petUnknown => 'Aún no descubierta';

  @override
  String get petAbout => 'Acerca de';

  @override
  String get petHabitat => 'Hábitat';

  @override
  String petElement(String name) {
    return 'Elemento: $name';
  }

  @override
  String get petStage1 => 'Cría';

  @override
  String get petStage2 => 'Joven';

  @override
  String get petStage3 => 'Adulta';

  @override
  String get bodyMale => 'Hombre';

  @override
  String get bodyFemale => 'Mujer';

  @override
  String get bodyChoice => 'Apariencia';

  @override
  String get bodyChange => 'Cambiar apariencia';

  @override
  String get portraitNoArmor => 'Sin armadura';

  @override
  String get combatPower => 'Poder de combate';

  @override
  String get combatTapHint => 'toca para detalles';

  @override
  String statCap(int cap) {
    return 'tope $cap%';
  }

  @override
  String get cEvasion => 'Evasión';

  @override
  String get cBlock => 'Bloqueo';

  @override
  String get cCritDmg => 'Daño crítico';

  @override
  String get cDmgReduction => 'Reducción de daño';

  @override
  String get cLifesteal => 'Robo de vida';

  @override
  String get cExecute => 'Ejecución';

  @override
  String get cMissingHp => 'Última resistencia';

  @override
  String get cPetAssist => 'Apoyo de mascota';

  @override
  String get xPower => 'La fuerza usada en los eventos: 2 × ataque + defensa + vida máx. / 10.';

  @override
  String get xAttack => 'Daño base de cada golpe antes de la defensa enemiga.';

  @override
  String get xDefense => 'Cada golpe que recibes pierde la mitad de tu defensa.';

  @override
  String xCrit(int cap) {
    return 'Probabilidad de que un golpe sea crítico. Como máximo $cap%.';
  }

  @override
  String get xCritDmg => 'Un crítico inflige este múltiplo de un golpe normal.';

  @override
  String get xHit => 'Tu probabilidad de acertar, menos la evasión enemiga (15–95%).';

  @override
  String xEvasion(int cap) {
    return 'Se resta de la probabilidad de acierto enemiga. Como máximo $cap%.';
  }

  @override
  String xBlock(int cap) {
    return 'Probabilidad de reducir a la mitad un golpe recibido. Como máximo $cap%.';
  }

  @override
  String xDmgReduction(int cap) {
    return 'Cada golpe recibido se reduce en esta parte. Como máximo $cap%.';
  }

  @override
  String get xLifesteal => 'Te curas esta parte del daño que infliges.';

  @override
  String xExecute(int pct) {
    return 'Daño extra contra enemigos con $pct% de vida o menos.';
  }

  @override
  String get xMissingHp => 'Hasta este daño extra a medida que baja tu vida.';

  @override
  String get xPetAssist => 'Tu mascota activa añade esta parte a cada golpe.';

  @override
  String get lookFace => 'Rostro';

  @override
  String get lookHairStyle => 'Peinado';

  @override
  String get lookHairColor => 'Color de pelo';

  @override
  String get lookEyeColor => 'Color de ojos';

  @override
  String get lookSave => 'Guardar';

  @override
  String get lookHumanOnly => 'Las opciones de rostro y pelo de esta raza llegarán pronto.';

  @override
  String get raceAbilities => 'Habilidades de raza';

  @override
  String get raceAbilitiesHint => 'todas pasivas';

  @override
  String get chooseClass => 'Elige tu clase';

  @override
  String raceAttrBonus(int pct, String name) {
    return '+$pct % $name';
  }

  @override
  String raceProfessionXp(int pct, String name) {
    return '+$pct % de XP extra en $name';
  }

  @override
  String get language => 'Idioma';

  @override
  String get report => 'Denunciar';

  @override
  String reportTitle(String name) {
    return 'Denunciar a $name';
  }

  @override
  String get reportWhy => '¿Qué ha pasado?';

  @override
  String get reportHarassment => 'Acoso o insultos';

  @override
  String get reportCheating => 'Trampas o bots';

  @override
  String get reportName => 'Nombre ofensivo';

  @override
  String get reportSpam => 'Spam o publicidad';

  @override
  String get reportScam => 'Estafa';

  @override
  String get reportOther => 'Otro motivo';

  @override
  String get reportNote => 'Detalles (opcional)';

  @override
  String get reportSent => 'Denuncia enviada. Gracias.';

  @override
  String get errReportSelf => 'No puedes denunciar a tu propio personaje.';

  @override
  String get nkReportReviewed => 'Tu denuncia se ha revisado';

  @override
  String get guide => 'Guía';

  @override
  String get guideSearch => 'Buscar en la guía';

  @override
  String get guideNoMatch => 'No se ha encontrado nada.';

  @override
  String get guideShowMe => 'Muéstramelo';

  @override
  String get next => 'Siguiente';

  @override
  String get back => 'Atrás';

  @override
  String get tutorialSkip => 'Saltar tutorial';

  @override
  String get tutorialStart => 'Empezar a jugar';

  @override
  String get tutorialReplay => 'Ver el tutorial otra vez';

  @override
  String get tutorialWelcomeTitle => 'Te damos la bienvenida a MeNa World';

  @override
  String get tutorialWelcomeBody =>
      'Tu personaje sigue trabajando con la aplicación cerrada. Dale una tarea, vuelve más tarde y recoge lo que haya ganado.';

  @override
  String get tutorialCharacterBody =>
      'Aquí están tus estadísticas, puntos de atributo y habilidades. Reparte los puntos que ganas en cada nivel y elige tus ascensos en los niveles 300 y 700.';

  @override
  String get tutorialGatherBody =>
      'Elige un recurso y una duración, y vuelve a recogerlo. Minería, pesca, tala, herboristería y caza suben de nivel por separado.';

  @override
  String get tutorialBagBody =>
      'Toca un objeto para equiparlo, usarlo, venderlo o publicarlo. La tabla comparativa muestra si supera a lo que llevas puesto.';

  @override
  String get tutorialAdventureBody =>
      'Caza monstruos, supera mazmorras, únete a incursiones y eventos, y viaja por el mapa. Los monstruos y mazmorras dependen de la región en la que estés.';

  @override
  String get tutorialWorkshopBody =>
      'Fabrica con tus profesiones, compra en la tienda y comercia con otros jugadores en el mercado.';

  @override
  String get tutorialTopBarTitle => 'La barra superior';

  @override
  String get tutorialTopBarBody =>
      'Misiones, notificaciones y social (chat, amigos, correo, gremio, intercambio) están arriba a la derecha. Un personaje hace una sola tarea a la vez: termínala o recógela antes de empezar otra.';

  @override
  String get tutorialDoneTitle => 'Ya estás listo';

  @override
  String get tutorialDoneBody => '¿Has olvidado dónde está algo? Abre la Guía al final de la pestaña Personaje.';

  @override
  String get guideCharacterTitle => 'Estadísticas, atributos y habilidades';

  @override
  String get guideCharacterBody =>
      'Pestaña Personaje. Reparte los puntos de atributo desde la tarjeta de progreso. Tus habilidades están en su propia sección plegable; de las que hacen lo mismo, se queda la más fuerte.';

  @override
  String get guidePromotionTitle => 'Ascensos';

  @override
  String get guidePromotionBody =>
      'Pestaña Personaje, tarjeta de progreso. En los niveles 300 y 700 la XP espera hasta que elijas uno de dos caminos. La XP ganada mientras tanto se paga al elegir.';

  @override
  String get guideGatherTitle => 'Recolección';

  @override
  String get guideGatherBody =>
      'Pestaña Recolección. Elige profesión, un recurso a tu alcance y 1, 3 o 6 horas. Los recursos pertenecen a regiones: viaja para conseguir mejores.';

  @override
  String get guideBagTitle => 'Bolsa y equipo';

  @override
  String get guideBagBody =>
      'Pestaña Bolsa. Toca un objeto para ver detalles y comparar. Mantén pulsado y arrastra para seleccionar varios y venderlos o destruirlos juntos. Con la bolsa llena no puedes empezar tareas.';

  @override
  String get guideBankTitle => 'Banco de la cuenta';

  @override
  String get guideBankBody =>
      'Pestaña Bolsa, botón Banco. Lo comparten todos tus personajes, así que sirve para pasarles objetos. No admite objetos ligados ni equipados.';

  @override
  String get guideHuntTitle => 'Caza';

  @override
  String get guideHuntBody =>
      'Pestaña Aventura, Caza. Elige un monstruo a tu alcance, una duración y las pociones y comida que llevarás. Tu personaje se retira solo si queda muy herido.';

  @override
  String get guideDungeonTitle => 'Mazmorras e incursiones';

  @override
  String get guideDungeonBody =>
      'Pestaña Aventura, Mazmorra. Las mazmorras son en solitario. Las incursiones piden un grupo de 2 a 5; usa Buscar grupo para emparejarte. Los jefes de incursión sueltan piezas de conjunto y huevos de mascota.';

  @override
  String get guideTravelTitle => 'Viajes y mapa';

  @override
  String get guideTravelBody =>
      'Pestaña Aventura, Viaje. Toca una región para ver sus monstruos, mazmorras, recursos y misiones, y luego Ir. Solo el destino exige un nivel.';

  @override
  String get guideEventsTitle => 'Eventos';

  @override
  String get guideEventsBody =>
      'Pestaña Aventura, Eventos. Alístate en jefes de mundo e invasiones y recoge la ficha diaria de temporada. Toca una tarjeta para ver cómo funciona y qué paga.';

  @override
  String get guideCraftTitle => 'Profesiones y fabricación';

  @override
  String get guideCraftBody =>
      'Pestaña Taller, Fabricación. Abre una profesión para ver las recetas. Los colores indican cuánta XP da aún una receta. Los pergaminos de encantamiento se usan desde la Bolsa.';

  @override
  String get guideVendorTitle => 'Tienda y recompra';

  @override
  String get guideVendorBody =>
      'Pestaña Taller, Tienda. Vende desde la Bolsa. Puedes recomprar lo vendido durante 2 horas.';

  @override
  String get guideMarketTitle => 'Mercado y subastas';

  @override
  String get guideMarketBody =>
      'Pestaña Taller, Mercado. Compra a precio fijo o puja en subastas. Publica objetos desde la Bolsa: precio fijo o subasta con puja inicial y compra inmediata opcional. Una puja en los últimos minutos alarga la subasta.';

  @override
  String get guideQuestsTitle => 'Misiones';

  @override
  String get guideQuestsBody =>
      'Icono de portapapeles arriba a la derecha. Misiones de historia, secundarias y diarias; la insignia cuenta las listas para entregar.';

  @override
  String get guideNotificationsTitle => 'Notificaciones y recompensas';

  @override
  String get guideNotificationsBody =>
      'Icono de campana arriba a la derecha. La pestaña Recompensas guarda lo recogido durante 30 días. En Ajustes eliges qué te llega.';

  @override
  String get guideSocialTitle => 'Chat, amigos, correo e intercambio';

  @override
  String get guideSocialBody =>
      'Globos de diálogo arriba a la derecha. Los intercambios los confirman ambas partes y pagan impuesto. Bloquear oculta a un jugador de tu chat, amigos y correo.';

  @override
  String get guidePartyTitle => 'Grupo';

  @override
  String get guidePartyBody =>
      'Pestaña Personaje, sección Grupo. Invita por nombre. Los miembros que hacen el mismo tipo de tarea ganan XP adicional.';

  @override
  String get guideGuildTitle => 'Gremio, tesorería y banco del gremio';

  @override
  String get guideGuildBody =>
      'Globos de diálogo arriba a la derecha, pestaña Gremio. Los miembros depositan oro y objetos; el rango decide quién retira y cuánto al día. Todo movimiento queda registrado.';

  @override
  String get guideRankingTitle => 'Clasificación';

  @override
  String get guideRankingBody =>
      'Pestaña Personaje, botón Clasificación. Tu puntuación viene del nivel, logros, maestrías, profesiones y mazmorras e incursiones superadas. Filtra por clase o por tabla.';

  @override
  String get guidePetsTitle => 'Mascotas';

  @override
  String get guidePetsBody =>
      'Pestaña Personaje, botón Mascotas. Los huevos salen de jefes de incursión de nivel alto. Incúbalos, súbelos de nivel cazando y activa uno por su bonificación y su ayuda en combate.';

  @override
  String get guideHelpersTitle => 'Ayudantes';

  @override
  String get guideHelpersBody =>
      'Pestaña Personaje, botón Ayudantes. Contrata un ayudante para una profesión. Trabaja por su cuenta, haga lo que haga tu personaje.';

  @override
  String get guideReportTitle => 'Denunciar a un jugador';

  @override
  String get guideReportBody =>
      'Mantén pulsada una línea del chat o usa el menú de un amigo. Elige un motivo; el equipo lo revisa y te avisa cuando lo haya visto.';

  @override
  String get guideAccountTitle => 'Cuenta, idioma y soporte';

  @override
  String get guideAccountBody =>
      'Lista de personajes, icono de cuenta arriba. Cambia el idioma o la contraseña, contacta con soporte, descarga tus datos o elimina la cuenta.';

  @override
  String get ranking => 'Clasificación';

  @override
  String get rankingScore => 'Puntuación';

  @override
  String get rankingLevel => 'Nivel';

  @override
  String get rankingAchievements => 'Logros';

  @override
  String get rankingMastery => 'Maestría';

  @override
  String get rankingProfessions => 'Profesiones';

  @override
  String get rankingClears => 'Superadas';

  @override
  String get rankingDungeons => 'Mazmorras';

  @override
  String get rankingRaids => 'Incursiones';

  @override
  String get rankingAllClasses => 'Todas las clases';

  @override
  String get rankingHowScore => 'Cómo se calcula la puntuación';

  @override
  String get rankingScoreNote => 'Cada mazmorra e incursión cuenta una vez. Se actualiza cada pocos minutos.';

  @override
  String get rankingEmpty => 'La tabla aún se está calculando.';

  @override
  String get guildVault => 'Tesorería y banco';

  @override
  String get guildVaultTreasury => 'Tesorería del gremio';

  @override
  String get guildVaultBank => 'Banco del gremio';

  @override
  String get guildVaultDepositGold => 'Depositar oro';

  @override
  String get guildVaultWithdrawGold => 'Retirar oro';

  @override
  String get guildVaultGoldAmount => 'Cantidad de oro';

  @override
  String get guildVaultNoLimit => 'Puedes retirar sin límite diario.';

  @override
  String guildVaultGoldToday(String used, String limit) {
    return 'Retirado hoy: $used de $limit';
  }

  @override
  String guildVaultItemsToday(int used, int limit) {
    return 'Objetos retirados hoy: $used de $limit';
  }

  @override
  String get guildVaultNoItems => 'Tu rango aún no puede usar el banco del gremio.';

  @override
  String get guildVaultLog => 'Movimientos';

  @override
  String get guildVaultLogEmpty => 'Aún no hay movimientos.';

  @override
  String guildVaultLogIn(String name, String what) {
    return '$name depositó $what';
  }

  @override
  String guildVaultLogOut(String name, String what) {
    return '$name retiró $what';
  }

  @override
  String get errGuildVaultNotEmpty => 'Vacía la tesorería y el banco antes de disolver el gremio.';

  @override
  String get errGuildVaultClosing => 'El gremio se está disolviendo: no se puede depositar.';

  @override
  String get errGuildVaultLimit => 'Eso supera tu límite diario.';

  @override
  String get errItemBound => 'Los objetos ligados se quedan con su personaje.';

  @override
  String get errItemEquipped => 'Quítate el objeto primero.';

  @override
  String get bid => 'Pujar';

  @override
  String get bidPlaced => 'Puja realizada.';

  @override
  String get marketFixed => 'Precio fijo';

  @override
  String get marketAuction => 'Subasta';

  @override
  String auctionNoBids(String price) {
    return 'Sin pujas, empieza en $price';
  }

  @override
  String auctionTopBid(String price, int count) {
    return 'Puja más alta $price ($count pujas)';
  }

  @override
  String get auctionYouLead => 'vas en cabeza';

  @override
  String get auctionStartBid => 'Puja inicial por todo el lote';

  @override
  String get auctionBuyoutOptional => 'Precio de compra inmediata (opcional)';

  @override
  String get auctionRules =>
      'Gana la puja más alta al acabar el tiempo. Una puja en los últimos 5 minutos alarga la subasta. No se puede cancelar una vez que alguien ha pujado.';

  @override
  String get auctionYouReceiveAtLeast => 'Recibes como mínimo';

  @override
  String auctionMinBid(String price) {
    return 'Puja mínima: $price';
  }

  @override
  String get auctionBidCopper => 'Tu puja en cobre';

  @override
  String get auctionHeldNote => 'El oro queda retenido y se devuelve en cuanto alguien te supere.';

  @override
  String auctionBuyoutFor(String price) {
    return 'Comprar ya: $price';
  }

  @override
  String get errBidTooLow => 'Esa puja está por debajo del mínimo.';

  @override
  String get errTopBidder => 'Ya tienes la puja más alta.';

  @override
  String get errAuctionBusy => 'Se está procesando otra puja. Inténtalo de nuevo en un momento.';

  @override
  String get errAuctionHasBids => 'No se puede cancelar una subasta con pujas.';

  @override
  String get nkMarketOutbid => 'Han superado tu puja';

  @override
  String get nkMarketWon => 'Has ganado una subasta';

  @override
  String get evMarketOutbid => 'Han superado tu puja; tu oro ha vuelto.';
}
