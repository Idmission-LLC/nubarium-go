package main

import (
    "encoding/json"
    "fmt"
    "io/fs"
    "math/rand"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

// Keys that contain personally identifiable or sensitive values
var piiKeyPattern = regexp.MustCompile(`(?i)^(nombre|calle|colonia|ciudad|cp|codigoPostal|numeroServicio|numeroMedidor|referencia|codigoBarras|codigoNumerico|codigoValidacion|qr|rmu2|multiplicador|tarifa|status|tipo|claveMensaje)$`)

// Keys that look like amounts that we want to normalize
var amountKeyPattern = regexp.MustCompile(`(?i)^(totalPagar|totalPagar2)$`)

// Keys that look like dates
var dateKeyPattern = regexp.MustCompile(`(?i)^(fecha|fechaLimitePago|periodoFacturado)$`)

func main() {
    // Default directory for fixtures
    targetDir := "testdata/responses"
    if len(os.Args) > 1 && os.Args[1] != "" {
        targetDir = os.Args[1]
    }

    if err := sanitizeDirectory(targetDir); err != nil {
        fmt.Fprintf(os.Stderr, "sanitize failed: %v\n", err)
        os.Exit(1)
    }
}

func sanitizeDirectory(dir string) error {
    return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        if !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
            return nil
        }
        if err := sanitizeFile(path); err != nil {
            return fmt.Errorf("%s: %w", path, err)
        }
        fmt.Printf("sanitized %s\n", path)
        return nil
    })
}

func sanitizeFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var v any
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }

    seed := int64(len(path))
    for i := 0; i < len(path); i++ {
        seed += int64(path[i])
    }
    rnd := rand.New(rand.NewSource(seed))

    sanitized := sanitizeValue("", v, filepath.Base(path), rnd)

    out, err := json.MarshalIndent(sanitized, "", "  ")
    if err != nil {
        return err
    }
    out = append(out, '\n')
    return os.WriteFile(path, out, 0644)
}

func sanitizeValue(key string, v any, context string, rnd *rand.Rand) any {
    switch val := v.(type) {
    case map[string]any:
        m := make(map[string]any, len(val))
        for k, vv := range val {
            m[k] = sanitizeValue(k, vv, context, rnd)
        }
        return m
    case []any:
        arr := make([]any, len(val))
        for i := range val {
            arr[i] = sanitizeValue(key, val[i], context, rnd)
        }
        return arr
    case string:
        if amountKeyPattern.MatchString(key) {
            return "100.00"
        }
        if dateKeyPattern.MatchString(key) {
            // Keep a realistic but anonymized date
            return fakeDate(rnd)
        }
        if piiKeyPattern.MatchString(key) {
            return placeholderForKey(key, context)
        }
        // Heuristic: emails, phones, long numeric strings
        if looksLikeEmail(val) {
            return "anon@example.com"
        }
        if looksLikePhone(val) {
            return "+520000000000"
        }
        if looksLikeBigNumber(val) {
            return strings.Repeat("X", len(val))
        }
        return val
    case float64:
        if amountKeyPattern.MatchString(key) {
            return 100.00
        }
        // For numeric identifiers that match PII keys, replace with deterministic number
        if piiKeyPattern.MatchString(key) {
            return float64(100000 + rnd.Intn(900000))
        }
        return val
    default:
        return val
    }
}

func placeholderForKey(key, context string) string {
    switch strings.ToLower(key) {
    case "nombre":
        return "ANON USER"
    case "calle":
        return "CALLE FALSA 123"
    case "colonia":
        return "COLONIA FALSA"
    case "ciudad":
        return "CIUDAD ANON"
    case "cp", "codigoPostal":
        return "00000"
    case "qr":
        return "QR-ANON-" + hashLike(context)
    case "codigoBarras":
        return "CB-ANON-" + hashLike(context)
    case "codigoNumerico", "codigoValidacion", "rmu2", "multiplicador", "referencia", "numeroServicio", "numeroMedidor":
        return "ANON-" + hashLike(context)
    case "tarifa":
        return "T-ANON"
    case "status":
        return "OK"
    case "tipo":
        return "ANON"
    case "claveMensaje":
        return "MENSAJE-ANON"
    default:
        return "ANON"
    }
}

func hashLike(s string) string {
    // Simple deterministic short token from input without importing crypto
    h := uint32(2166136261)
    for i := 0; i < len(s); i++ {
        h ^= uint32(s[i])
        h *= 16777619
    }
    return fmt.Sprintf("%08x", h)
}

func looksLikeEmail(s string) bool {
    return strings.Contains(s, "@") && strings.Contains(s, ".")
}

func looksLikePhone(s string) bool {
    digits := 0
    for i := 0; i < len(s); i++ {
        if s[i] >= '0' && s[i] <= '9' {
            digits++
        }
    }
    return digits >= 10 && digits <= 15
}

func looksLikeBigNumber(s string) bool {
    if len(s) < 8 {
        return false
    }
    for i := 0; i < len(s); i++ {
        if s[i] < '0' || s[i] > '9' {
            return false
        }
    }
    return true
}

func fakeDate(rnd *rand.Rand) string {
    // Generate a date within the last 2 years in format YYYY-MM-DD
    now := time.Now()
    days := rnd.Intn(365 * 2)
    d := now.AddDate(0, 0, -days)
    return d.Format("2006-01-02")
}


