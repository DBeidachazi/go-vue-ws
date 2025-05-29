<template>
  <h1>{{ msg }}</h1>

  <div class="voting-container">
    <h2>请投票选择您喜欢的选项：</h2>
    <div class="voting-buttons">
      <button
          v-for="item in voteData"
          :key="item.id"
          @click="vote(item.id)"
          :class="{ 'voted': hasVoted && votedOptionId === item.id }"
      >
        {{ item.name }} ({{ item.value }})
      </button>
    </div>
    <p v-if="hasVoted" class="vote-message">您已经投票给"{{ getOptionNameById(votedOptionId) }}"</p>
    <p v-if="loading" class="loading-message">正在加载数据...</p>
    <p v-if="error" class="error-message">{{ error }}</p>
  </div>

  <Charts :data="voteData"/>
</template>

<script setup lang="ts">
import Charts from "./Charts.vue";
import {onMounted, ref, reactive} from 'vue'
import service from "../util/axios.ts";
import { WebSocketClient } from "../util/ws";
import {getSessionId} from '../util/auth';

defineProps<{ msg: string }>()

interface VoteOption {
  id: number;
  name: string;
  value: number;
}

let voteData = reactive<VoteOption[]>([])
const hasVoted = ref(false);
const votedOptionId = ref(0);
const loading = ref(false);
const error = ref('');

// 创建WebSocket客户端
const wsClient = new WebSocketClient('ws://localhost:3000/ws/poll');

const handleWSMessage = (data: any) => {
  console.log('收到WebSocket消息:', data);
  if (data && data.polls) {
    // 更新投票数据
    updateVoteData(data.polls);
  }
};

const handleWSError = (error: any) => {
  console.error('WebSocket连接错误:', error);
  error.value = '实时数据连接失败，正在使用本地数据';
};

// ws后处理
const updateVoteData = (polls: any[]) => {
  if (!polls || !Array.isArray(polls)) return;

  voteData.length = 0;

  polls.forEach((item) => {
    voteData.push({
      id: item.option_id,
      name: item.option_description,
      value: item.vote_count
    });
  });
};

// 启动处理
const fetchVoteData = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await service.get('/api/poll');
    voteData.length = 0;

    if (response.data && response.data.polls) {
      response.data.polls.forEach((item: any) => {
        voteData.push({
          id: item.option_id,
          name: item.option_description,
          value: item.vote_count
        });
      });
    } else {
      throw new Error('服务器返回的数据格式不正确');
    }
  } catch (err: any) {
    console.error('获取投票数据失败:', err);
    error.value = `获取数据失败: ${err.message || '网络错误'}`;

    // default
    if (voteData.length === 0) {
      const defaultData: VoteOption[] = [
        {id: 1, name: '选项A', value: 120},
        {id: 2, name: '选项B', value: 200},
        {id: 3, name: '选项C', value: 150},
        {id: 4, name: '选项D', value: 80}
      ];
      defaultData.forEach(item => voteData.push(item));
    }
  } finally {
    loading.value = false;
  }
};

const vote = (id: number) => {
  if (hasVoted.value) {
    alert('您已经投过票了！');
    return;
  }

  const option = voteData.find(item => item.id === id);
  if (option) {
    option.value += 1;

    hasVoted.value = true;
    votedOptionId.value = id;

    service.post('/api/poll/vote', { user_uuid: getSessionId(), vote_option: id })
      .then(() => {
        console.log(`投票成功，选项ID: ${id}`);
      })
      .catch(error => {
        console.error('投票失败:', error);
      });
    sessionStorage.setItem('voted_option', id.toString());
  }
};

const getOptionNameById = (id: number): string => {
  const option = voteData.find(item => item.id === id);
  return option ? option.name : '';
};

onMounted(() => {
  // todo 刷新清除测试数据
  sessionStorage.removeItem('browser_session_id');
  sessionStorage.removeItem('voted_option');

  const sessionId = getSessionId();
  console.log(`当前会话ID: ${sessionId}`);

  fetchVoteData();

  const votedOption = sessionStorage.getItem('voted_option');
  if (votedOption) {
    hasVoted.value = true;
    votedOptionId.value = parseInt(votedOption);
  }

  // WebSocket事件监听
  wsClient.on('message', handleWSMessage);
  wsClient.on('error', handleWSError);
  wsClient.connect();
});

</script>

<style scoped>
.voting-container {
  margin: 20px 0;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 5px;
}

.voting-buttons {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin: 15px 0;
}

button {
  padding: 10px 20px;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s;
}

button:hover {
  background-color: #45a049;
  transform: translateY(-2px);
}

button.voted {
  background-color: #2196F3;
  font-weight: bold;
}

.vote-message {
  color: #2196F3;
  font-weight: bold;
  text-align: center;
}

.loading-message {
  color: #ff9800;
  text-align: center;
}

.error-message {
  color: #f44336;
  text-align: center;
}
</style>