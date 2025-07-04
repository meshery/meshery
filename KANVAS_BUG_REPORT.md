# Kanvas-Snapshot Buffer Allocation Bug Report

## Issue Summary
The Kanvas-Snapshot action consistently fails with a buffer allocation error across multiple versions (v0.2.37, v0.2.39), preventing CI workflows from completing successfully.

## Error Details
```
RangeError: The argument 'size' is invalid. Received -6667200
RangeError [ERR_INVALID_ARG_VALUE]: The argument 'size' is invalid. Received -6667200
    at Function.allocUnsafe (node:buffer:376:3)
    at s.cropQuiet (<embedded>:2186:51151)
    at W.<computed> [as crop] (<embedded>:1981:692110)
```

## Technical Analysis

### Environment
- **Node.js**: v20.19.1
- **Cypress**: 12.13.0
- **Browser**: Chrome 138 (headless)
- **Viewport**: 1920x1080
- **Screenshot dimensions**: 5760x2823 pixels

### Root Cause
The error occurs in the image cropping function where a negative buffer size is calculated:
- **Buffer size**: -6667200 bytes (negative value)
- **Expected calculation**: width × height × bytes_per_pixel
- **Likely issue**: Integer overflow or incorrect crop coordinate calculation

### Mathematical Analysis
```
Screenshot: 5760 × 2823 = 16,273,920 pixels
Expected buffer: 16,273,920 × 4 bytes = 65,095,680 bytes (~65MB)
Actual error: -6667200 bytes (negative)
```

This suggests either:
1. **Crop coordinates are negative** (invalid crop region)
2. **Integer overflow** in 32-bit signed calculation
3. **Buffer size arithmetic error** in the image processing library

## Impact
- ❌ **CI workflows fail completely**
- ❌ **Pull requests cannot be merged**
- ❌ **Screenshots are not generated**
- ❌ **Infrastructure validation blocked**

## Attempted Fixes
1. ✅ Updated to latest version (v0.2.39)
2. ✅ Increased Node.js heap size (`--max-old-space-size=8192`)
3. ✅ Added Cypress cache configuration
4. ✅ Added `continue-on-error: true` to prevent workflow failures

## Current Status
- **Workaround**: `continue-on-error: true` prevents complete CI failure
- **Long-term**: Requires upstream fix in Kanvas-Snapshot action
- **Monitoring**: Issue persists across multiple versions

## Recommendations

### Immediate (for users)
1. Add `continue-on-error: true` to Kanvas-Snapshot steps
2. Use Node.js memory optimization settings
3. Monitor for upstream fixes

### Long-term (for maintainers)
1. **Fix crop calculation logic** in image processing
2. **Add input validation** for crop coordinates
3. **Implement error bounds checking** for buffer allocation
4. **Add fallback mechanisms** for large screenshots
5. **Consider viewport size limits** to prevent large screenshots

## Files Modified
- `.github/workflows/kanvas.yml` - Added error handling and memory config
- `docs/pages/extensions/kanvas-snapshot/index.md` - Updated to v0.2.39
- `fix_kanvas_tests.sh` - Diagnostic script
- `kanvas_advanced_diagnostic.sh` - Detailed analysis tool

## Related Issues
- Buffer allocation errors in image processing libraries
- Screenshot dimension calculation bugs
- Memory management in headless browser automation

---
**Status**: Ongoing upstream issue  
**Severity**: High (blocks CI)  
**Workaround**: Available (continue-on-error)  
**Fix Required**: Upstream (Kanvas-Snapshot action)
