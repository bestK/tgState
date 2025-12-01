module.exports = {
  apps: [{
    name: 'tgstate',
    script: './tgstate',
    args: '-port 8088 -mode p',
    cwd: './',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '1G',
    env: {
      NODE_ENV: 'production'
    },
    error_file: './logs/tgstate-error.log',
    out_file: './logs/tgstate-out.log',
    log_file: './logs/tgstate-combined.log',
    time: true,
    merge_logs: true,
    log_date_format: 'YYYY-MM-DD HH:mm:ss Z'
  }]
};
