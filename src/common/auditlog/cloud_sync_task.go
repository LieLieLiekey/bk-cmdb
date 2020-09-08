/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auditlog

import (
	"fmt"

	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type syncTaskAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of cloud sync task, if data is nil, will auto get data by taskID.
func (h *syncTaskAuditLog) GenerateAuditLog(kit *rest.Kit, action metadata.ActionType, taskID int64, OperateFrom metadata.OperateFromType,
	data *metadata.CloudSyncTask, updateFields map[string]interface{}) (*metadata.AuditLog, error) {
	if data == nil {
		// get data by taskID.
		option := metadata.SearchCloudOption{
			Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID},
		}

		res, err := h.clientSet.Cloud().SearchSyncTask(kit.Ctx, kit.Header, &option)
		if err != nil {
			blog.Errorf("generate audit log of cloud sync task, failed to read cloud sync task, err: %v, rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}
		if len(res.Info) <= 0 {
			blog.Errorf("generate audit log of cloud sync task failed, not find cloud sync task, rid: %s",
				kit.Rid)
			return nil, fmt.Errorf("generate audit log of cloud sync task failed, not find cloud sync task")
		}

		data = &res.Info[0]
	}

	taskName := data.TaskName

	var basicDetail *metadata.BasicContent
	switch action {
	case metadata.AuditCreate:
		basicDetail = &metadata.BasicContent{
			CurData: data.ToMapStr(),
		}
	case metadata.AuditDelete:
		basicDetail = &metadata.BasicContent{
			PreData: data.ToMapStr(),
		}
	case metadata.AuditUpdate:
		basicDetail = &metadata.BasicContent{
			PreData:      data.ToMapStr(),
			UpdateFields: updateFields,
		}
	}

	var auditLog = &metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudSyncTaskRes,
		Action:       action,
		ResourceID:   taskID,
		ResourceName: taskName,
		OperateFrom:  OperateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: basicDetail,
		},
	}

	return auditLog, nil
}

func NewSyncTaskAuditLog(clientSet coreservice.CoreServiceClientInterface) *syncTaskAuditLog {
	return &syncTaskAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
