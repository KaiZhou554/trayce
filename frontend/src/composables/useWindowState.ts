import { ref } from 'vue'
import {
  WindowIsMaximised,
  WindowMaximise,
  WindowUnmaximise,
} from '../../wailsjs/runtime/runtime'

/**
 * Reactive window-state synchronisation.
 *
 * Because the user can maximise/restore the window via OS gestures
 * (double-click title bar, Win+↑/↓, Snap Layout, drag-to-restore, …)
 * the frontend has no native event to listen to.
 *
 * This composable polls `WindowIsMaximised()` on a fixed interval and
 * exposes a reactive `isMaximized` ref that stays in sync with the OS.
 *
 * Future extensions (isFullscreen, isMinimized, isFocused) can be added
 * by extending `syncWindowState()` and adding extra refs.
 */
export function useWindowState() {
  // ── Reactive state ──
  const isMaximized = ref(false)

  // ── Polling timer ──
  let timer: ReturnType<typeof setInterval> | null = null

  // ── Core: query the OS and update reactive state ──
  async function syncWindowState(): Promise<boolean> {
    const maximized = await WindowIsMaximised()
    if (maximized !== isMaximized.value) {
      isMaximized.value = maximized
    }
    return maximized
  }

  // ── Lifecycle: start / stop polling ──
  function startWindowStateSync(intervalMs = 500) {
    stopWindowStateSync()
    syncWindowState() // immediate first sync
    timer = setInterval(syncWindowState, intervalMs)
  }

  function stopWindowStateSync() {
    if (timer !== null) {
      clearInterval(timer)
      timer = null
    }
  }

  // ── Actions: maximise / restore, then re-sync from OS ──
  async function maximize() {
    await WindowMaximise()
    await syncWindowState()
  }

  async function unmaximize() {
    await WindowUnmaximise()
    await syncWindowState()
  }

  return {
    isMaximized,
    syncWindowState,
    startWindowStateSync,
    stopWindowStateSync,
    maximize,
    unmaximize,
  }
}
