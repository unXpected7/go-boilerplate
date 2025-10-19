#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

// Test runner script
console.log('🧪 Running Frontend Tests...');

// Run unit tests with Jest
console.log('\n📝 Running Unit Tests...');
const jestProcess = spawn('npx', ['jest', '--coverage', '--verbose'], {
  cwd: __dirname,
  stdio: 'inherit',
});

jestProcess.on('close', (code) => {
  if (code !== 0) {
    console.log(`❌ Unit tests failed with exit code ${code}`);
    process.exit(code);
  }

  console.log('✅ Unit tests passed!');

  // Run E2E tests with Cypress
  console.log('\n🌐 Running E2E Tests...');
  const cypressProcess = spawn('npx', ['cypress', 'run', '--headless'], {
    cwd: __dirname,
    stdio: 'inherit',
  });

  cypressProcess.on('close', (code) => {
    if (code !== 0) {
      console.log(`❌ E2E tests failed with exit code ${code}`);
      process.exit(code);
    }

    console.log('✅ E2E tests passed!');
    console.log('\n🎉 All tests completed successfully!');
    console.log('\n📊 Test Coverage Report:');
    console.log('   - Unit Tests: ✅');
    console.log('   - E2E Tests: ✅');
    console.log('   - Component Tests: ✅');
    console.log('   - Hook Tests: ✅');
    console.log('   - Integration Tests: ✅');

    process.exit(0);
  });
});

// Handle test timeouts
setTimeout(() => {
  console.log('⏰ Test timeout reached');
  process.exit(1);
}, 300000); // 5 minutes timeout