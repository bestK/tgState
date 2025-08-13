import axios from "axios";

export interface UploadResponse {
  code: number;
  message: string;
  imgUrl?: string;
  proxyUrl?: string;
  shortUrl?: string;
  shortFileUrl?: string;
  name?: string;
  chunkId?: string;
}

export interface MergeRequest {
  uploadId: string;
  fileName: string;
  chunkIds: string[];
  fileSize: number;
}

export interface FileRecord {
  fileId: string;
  filename: string;
  ip: string;
  time: string;
  userFingerprint?: string;
  shared: boolean;
}

export interface PaginationInfo {
  page: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}

export interface HistoryResponse {
  code: number;
  message: string;
  data?: {
    files: FileRecord[];
    pagination: PaginationInfo;
  };
}

const api = axios.create({
  baseURL: "/",
  timeout: 30000,
});

export const uploadFile = async (
  file: File,
  userFingerprint?: string,
  shared?: boolean,
  onProgress?: (progress: number) => void
): Promise<UploadResponse> => {
  const formData = new FormData();
  formData.append("file", file);
  if (userFingerprint) {
    formData.append("userFingerprint", userFingerprint);
  }
  if (shared !== undefined) {
    formData.append("shared", shared.toString());
  }

  const response = await api.post<UploadResponse>("/api", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total && onProgress) {
        const progress = (progressEvent.loaded / progressEvent.total) * 100;
        onProgress(progress);
      }
    },
  });

  return response.data;
};

export const uploadChunk = async (
  chunk: Blob,
  chunkIndex: number,
  uploadId: string,
  fileName: string,
  userFingerprint?: string,
  onProgress?: (progress: number) => void
): Promise<UploadResponse> => {
  const formData = new FormData();
  formData.append("file", chunk, `${fileName}.chunk.${chunkIndex}`);
  formData.append("chunkIndex", chunkIndex.toString());
  formData.append("uploadId", uploadId);
  formData.append("fileName", fileName);
  if (userFingerprint) {
    formData.append("userFingerprint", userFingerprint);
  }

  const response = await api.post<UploadResponse>("/api/chunk", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
    onUploadProgress: (progressEvent) => {
      if (progressEvent.total && onProgress) {
        const progress = (progressEvent.loaded / progressEvent.total) * 100;
        onProgress(progress);
      }
    },
  });

  return response.data;
};

export const mergeChunks = async (
  request: MergeRequest & { userFingerprint?: string; shared?: boolean }
): Promise<UploadResponse> => {
  const response = await api.post<UploadResponse>("/api/merge", request, {
    headers: {
      "Content-Type": "application/json",
    },
  });

  return response.data;
};

// 获取用户历史文件
export const getUserHistory = async (
  userFingerprint: string,
  page = 1,
  pageSize = 20
): Promise<HistoryResponse> => {
  const response = await api.get<HistoryResponse>(
    `/api/history?fingerprint=${encodeURIComponent(userFingerprint)}&page=${page}&pageSize=${pageSize}`
  );
  return response.data;
};

// 获取广场文件
export const getPlazaFiles = async (
  page = 1,
  pageSize = 20
): Promise<HistoryResponse> => {
  const response = await api.get<HistoryResponse>(
    `/api/plaza?page=${page}&pageSize=${pageSize}`
  );
  return response.data;
};
