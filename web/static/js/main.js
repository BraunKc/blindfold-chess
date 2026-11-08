class MoveEntry {
  constructor(input, cells, moveBtn) {
    this.input = input
    this.cells = cells
    this.moveBtn = moveBtn

    this.#addListeners()
  }

  onFocus(e) {
    const len = e.target.value.length
    e.target.setSelectionRange(len, len)
  }

  onInput(e) {
    this.#validateInput(e)

    const cell = e.target
    if (cell.value !== '') {
      const nextCell = this.#getNextCell(cell)
      if (nextCell) nextCell.focus()
    }
  }

  onKeydown(e) {
    if (e.key === 'Enter') this.#onEnter(e)
    else if (e.key === 'Backspace') this.#onBackspace(e)
    else if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') this.#onArrows(e)
  }

  clear() {
    this.cells.forEach((el) => {
      el.value = ''
    })
  }

  async makeMove(notation) {
    if (!notation) return

    const res = await fetch('/api/move', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        notation: notation,
      }),
    })

    const data = await res.json()
    console.log(data)
    if (!res.ok) {
      throw new Error(data.error || `HTTP ${res.status}`)
    }

    return data
  }

  #getNextCell(cell) {
    const cellIndex = Number(cell.dataset.index)
    return cellIndex < this.cells.length - 1
      ? this.cells[cellIndex + 1]
      : undefined
  }

  #getPrevCell(cell) {
    const cellIndex = Number(cell.dataset.index)
    return cellIndex > 0 ? this.cells[cellIndex - 1] : undefined
  }

  #getNotation() {
    let notation = ''
    this.cells.forEach((el) => {
      notation += el.value.toLowerCase()
    })

    return notation.length === 4 ? notation : null
  }

  async #onEnter(e) {
    e.preventDefault()

    const cell = e.target
    if (cell.value === '') return

    const notation = this.#getNotation()
    if (notation) {
      await this.makeMove(notation)
      return
    }

    const nextCell = this.#getNextCell(cell)
    if (nextCell) nextCell.focus()
  }

  #onBackspace(e) {
    e.preventDefault()

    const cell = e.target
    if (cell.value !== '') {
      cell.value = ''
      return
    }

    const prevCell = this.#getPrevCell(cell)
    if (!prevCell) return

    prevCell.focus()
    prevCell.value = ''
  }

  #onArrows(e) {
    e.preventDefault()

    const cellIndex = Number(e.target.dataset.index)
    const next = e.key === 'ArrowLeft' ? cellIndex - 1 : cellIndex + 1

    if (next >= 0 && next < this.cells.length) {
      this.cells[next].focus()
    }
  }

  #validateInput(e) {
    const cell = e.target
    const cellIndex = Number(cell.dataset.index)
    let value = cell.value.slice(0, 1)

    const isLetterCell = cellIndex % 2 === 0
    const valid = isLetterCell ? /^[a-h]$/i.test(value) : /^[1-8]$/.test(value)

    if (value && !valid) {
      cell.value = ''
      return
    }

    cell.value = value.toLowerCase()
  }

  #addListeners() {
    this.input.addEventListener('focusin', (e) => this.onFocus(e))
    this.input.addEventListener('input', (e) => this.onInput(e))
    this.input.addEventListener('keydown', (e) => this.onKeydown(e))
    this.moveBtn.addEventListener('click', () =>
      this.makeMove(this.#getNotation()),
    )
  }
}

const inputEl = document.getElementById('moveInput')
const cellsEls = document.querySelectorAll('.move-cell')
const moveBtnEl = document.getElementById('moveBtn')

const moveEntry = new MoveEntry(inputEl, cellsEls, moveBtnEl)
