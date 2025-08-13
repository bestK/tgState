import FingerprintJS from "@fingerprintjs/fingerprintjs";

// FingerprintJS 实例缓存
let fpInstance: any = null;

// 初始化 FingerprintJS
async function initFingerprint() {
  if (!fpInstance) {
    fpInstance = await FingerprintJS.load();
  }
  return fpInstance;
}

// 生成指纹哈希
export async function generateFingerprintHash(): Promise<string> {
  try {
    const fp = await initFingerprint();
    const result = await fp.get();
    return result.visitorId;
  } catch (error) {
    console.warn("FingerprintJS failed, falling back to simple hash:", error);
    // 降级方案：使用简单的浏览器信息生成哈希
    return await generateFallbackHash();
  }
}

// 降级方案：简单的指纹生成
async function generateFallbackHash(): Promise<string> {
  const data = {
    userAgent: navigator.userAgent,
    language: navigator.language,
    screenResolution: `${screen.width}x${screen.height}`,
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    colorDepth: screen.colorDepth,
    pixelRatio: window.devicePixelRatio,
    cookieEnabled: navigator.cookieEnabled,
  };

  const dataString = JSON.stringify(data);
  const encoder = new TextEncoder();
  const dataBuffer = encoder.encode(dataString);
  const hashBuffer = await crypto.subtle.digest("SHA-256", dataBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");

  return hashHex;
}

// 获取或生成用户指纹ID
export async function getUserFingerprintId(): Promise<string> {
  const storageKey = "tgstate_user_fingerprint";

  // 尝试从localStorage获取已存储的指纹
  let fingerprintId = localStorage.getItem(storageKey);

  if (!fingerprintId) {
    // 生成新的指纹ID
    fingerprintId = await generateFingerprintHash();
    localStorage.setItem(storageKey, fingerprintId);
  }

  return fingerprintId;
}
