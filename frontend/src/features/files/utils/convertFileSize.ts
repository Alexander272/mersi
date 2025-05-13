/**
 * @param size размер файла в байтах
 * @param digits кол-во знаком после запятой
 * @returns размер файла в мегабайтах
 */
export const convertFileSize = (size: number, digits: number) => (size / Math.pow(1024, 2)).toFixed(digits)
