const api = new API()

const turnEl = document.getElementById('turn')
const lastMovesEl = document.getElementById('lastMoves')
const stateView = new StateView(turnEl, lastMovesEl)

const notificationEl = document.getElementById('notification')
const notifier = new Notifier(notificationEl)

const promotionModalEl = document.getElementById('promotionModal')
const promotionBackdropEl = document.getElementById('promotionBackdrop')
const promotionModal = new PromotionModal(promotionModalEl, promotionBackdropEl)

const game = new Game(api, stateView, promotionModal, notifier)
try {
  await game.init()
} catch (err) {
  notifier.showError(err.message)
}

const inputEl = document.getElementById('moveInput')
const cellsEls = document.querySelectorAll('.move-cell')
const moveBtnEl = document.getElementById('moveBtn')
new MoveEntry(game, inputEl, cellsEls, moveBtnEl)
